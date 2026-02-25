# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

Phabricator — a PHP web application suite for software development (code review, repository browsing, task tracking, wiki, etc.). No longer actively maintained upstream, but this fork has ongoing PHP 8 compatibility fixes.

## Environment & Dependencies

**This checkout is source-code-only** — no PHP, Arcanist, or MySQL are installed. All 38 `bin/` commands are `#!/usr/bin/env php` and cannot execute. Node.js v22 is available (Aphlict server could theoretically run standalone).

To make it functional you need: PHP 8.x, Arcanist (sibling `../arcanist/` dir), libphutil (sibling `../libphutil/` dir), MySQL/MariaDB, and a web server.

Required PHP extensions: hash, json, openssl, mbstring, ctype, curl, mysqli. Optional: gd, fileinfo, zip, pcntl, posix, OPcache.

All third-party dependencies are vendored in `externals/` (13 libraries including stripe-php, phpmailer, xhprof, JsShrink). No Composer/npm needed.

## Development Commands

```bash
arc lint                              # Lint changed files
arc lint --everything                 # Lint entire repo
arc lint path/to/file.php             # Lint specific file
arc unit                              # Run tests for changed files
arc unit src/path/to/__tests__/       # Run tests in a directory
arc unit path/to/SomeTestCase.php     # Run a single test file
arc liberate src/                     # Regenerate library map (after adding/moving/removing classes)
bin/storage upgrade                   # Apply pending DB migrations
bin/celerity map                      # Rebuild static resource map (after changing JS/CSS)
bin/phd start|stop|restart|status     # Daemon management
bin/phd debug <DaemonClass>           # Run single daemon in foreground
```

## Architecture

### Request Lifecycle

```
webroot/index.php
  → PhabricatorStartup::loadCoreLibraries()
  → AphrontApplicationConfiguration::runHTTPRequest($sink)
    → Site resolution (PhabricatorPlatformSite/ResourceSite/ShortSite by host priority)
    → Route matching: nested regex arrays from PhabricatorApplication::getRoutes()
    → PhabricatorController::willBeginExecution()
      13-step enforcement chain: session → CSRF → disabled user → partial session
      → Legalpad → MFA → login required → email verification → spaces → app policies → admin
    → Controller::handleRequest($request)
    → Response production (AphrontResponseProducerInterface chains resolved iteratively)
    → AphrontHTTPSink::writeResponse() with CSP headers
```

Exception handling: priority-ordered handler chain (Policy 320k → HighSecurity 310k → RateLimit 300k → Ajax 110k → Conduit 100k → Default 900k). Multi-level fallback down to bare `PhabricatorStartup::didEncounterFatalException()`.

### Source Layout

- **`src/applications/`** — 76 self-contained application modules
- **`src/aphront/`** — Web framework: controllers, requests, 14 response types, HTTP parameter types, routing, write guard
- **`src/infrastructure/`** — Lisk ORM, database clustering, query system, edges, Ferret search, caching, markup/Remarkup, daemons, workers
- **`src/view/`** — Shared UI components (AphrontView)
- **`webroot/rsrc/`** — Frontend: JS (Javelin framework + behaviors), CSS (PHUI components), images, fonts
- **`scripts/`** — Implementations for `bin/` CLI entry points
- **`resources/sql/autopatches/`** — 1100+ database migration files
- **`resources/celerity/`** — Generated resource map + package definitions
- **`externals/`** — Vendored third-party libraries
- **`support/`** — Bootstrap, linting config, Aphlict Node.js websocket server

### Application Module Structure

Each app in `src/applications/{app}/` follows a standard layout:

| Directory | Base Class | Purpose |
|-----------|-----------|---------|
| `application/` | `PhabricatorApplication` | Routes, policies, menu items |
| `controller/` | `PhabricatorController` | Request handlers (`handleRequest()`) |
| `storage/` | `PhabricatorLiskDAO` (via app DAO base) | ORM objects + `SchemaSpec` |
| `query/` | `PhabricatorCursorPagedPolicyAwareQuery` | Policy-aware data loading |
| `editor/` | `PhabricatorApplicationTransactionEditor` | Mutation orchestration |
| `xaction/` | `PhabricatorModularTransactionType` | Per-field transaction types |
| `conduit/` | `ConduitAPIMethod` | API endpoints |
| `view/` | `AphrontView` | Application-specific views |
| `herald/` | `HeraldAdapter` | Rule adapters/actions |
| `search/` | `PhabricatorSearchEngineExtension` | Search index definitions |

### Lisk ORM

Object-authoritative ORM — `protected` properties auto-map to DB columns via Reflection. `private` properties are NOT persisted.

```php
class ManiphestTask extends ManiphestDAO {
  protected $title;        // → column `title`, auto getter/setter via __call()
  protected $status;       // → column `status`
  protected $ownerPHID;    // → column `ownerPHID`
  // Built-in: $id, $phid, $dateCreated, $dateModified
}
```

Key `getConfiguration()` options: `CONFIG_AUX_PHID` (auto-generate PHID), `CONFIG_SERIALIZATION` (JSON/PHP serialize columns), `CONFIG_COLUMN_SCHEMA` / `CONFIG_KEY_SCHEMA` (declarative schema).

Save lifecycle: `willSaveObject()` → `willWriteData()` (serialize) → SQL INSERT/UPDATE → `didWriteData()`. Column types inferred from names (`*PHID` → phid, `*ID` → id, `viewPolicy` → policy, `dateCreated` → epoch).

Table name derived by stripping `Phabricator` prefix: `ManiphestTask` → `maniphest_task` in DB `phabricator_maniphest`.

### Query System

```
PhabricatorQuery → PhabricatorOffsetPagedQuery → PhabricatorPolicyAwareQuery
  → PhabricatorCursorPagedPolicyAwareQuery → ConcreteQuery
```

Execution loop: `loadPage()` → `willFilterPage()` (pre-policy) → `PhabricatorPolicyFilter::apply()` → `didFilterPage()` (post-policy) → advance cursor → repeat until enough results.

**Overheating**: stops after examining 10x the needed rows to prevent runaway queries when viewer can see few objects. **Cursor paging**: multi-column cursor support for compound orderings. **Query workspace**: per-viewer cache preventing redundant loads during nested policy checks.

### Policy System

Objects implement `PhabricatorPolicyInterface` (3 methods: `getCapabilities()`, `getPolicy($cap)`, `hasAutomaticCapability($cap, $viewer)`). Four capabilities: `CAN_VIEW`, `CAN_EDIT`, `CAN_JOIN`, `CAN_INTERACT`.

Policy values: global constants (`public`/`users`/`admin`/`no-one`), user PHIDs, project PHIDs (membership check), or custom policy PHIDs (`PHID-PLCY-xxx` with rule-based logic stored in DB).

`PhabricatorPolicyFilter` evaluation order: omnipotent bypass → spaces check → `hasAutomaticCapability()` → policy switch. Extended policies (`PhabricatorExtendedPolicyInterface`) enable cross-object dependencies with depth-32 cycle protection.

### Transaction System

All mutations flow through `PhabricatorApplicationTransactionEditor::applyTransactions()`:

```
Lock & reload → expand → validate → compute old/new → capability checks → filter no-ops
  → BEGIN TX: applyInternalEffects → save object → save xaction rows → applyExternalEffects → COMMIT
  → Herald rules (post-commit, via sub-editor with omnipotent user)
  → Queue publish worker (mail, feed, search index, webhooks — async)
```

Modern transaction types extend `PhabricatorModularTransactionType` with `TRANSACTIONTYPE` constant and implement: `generateOldValue()`, `generateNewValue()`, `validateTransactions()`, `applyInternalEffects()`, `applyExternalEffects()`, display methods (`getTitle()`, `getIcon()`, `getColor()`).

### Frontend: Celerity + Javelin

**Celerity** manages static resources (JS/CSS). Files declare dependencies via docblock `@provides`/`@requires`. Resource map at `resources/celerity/map.php` (rebuilt via `bin/celerity map`).

PHP views declare dependencies:
```php
require_celerity_resource('phui-button-css');          // CSS/JS symbol
Javelin::initBehavior('my-feature', $config_array);    // JS behavior + config
```

**Javelin** is the custom JS framework: `JX.install()` (class system), `JX.Stratcom` (event delegation via `data-sigil`/`data-meta` attributes), `JX.behavior()` (server→client glue pattern), `JX.Workflow` (form/dialog workflows).

CSS theming via `{$variableName}` syntax resolved at serve-time by `CelerityResourceTransformer`. Five postprocessors: default, darkmode, highcontrast, largefont, redgreen. No Sass/Less — PHP handles variable substitution.

Packages defined in `resources/celerity/packages.php` (core.pkg.js, core.pkg.css, per-app packages).

### Aphlict (Real-time Notifications)

Node.js WebSocket server (`support/aphlict/server/`). Zero npm dependencies — uses only Node built-in modules.

```
PHP backend → HTTP POST → AphlictAdminServer → transmit → AphlictClientServer (WebSocket) → browser
```

Cross-tab leader election via `JX.Leader` (localStorage-based) — only one tab holds the WebSocket connection. Message replay (4096 messages / 60s history) for reconnection.

### Daemons & Background Processing

`bin/phd start` launches 4 daemons:
1. **PhabricatorTaskmasterDaemon** — processes worker task queue (autoscale pool)
2. **PhabricatorTriggerDaemon** — scheduled/recurring events + garbage collection
3. **PhabricatorRepositoryPullLocalDaemon** — keeps repo working copies updated
4. **PhabricatorFactDaemon** — computes analytics datapoints

**Worker task queue**: MySQL-backed, lease-based (`PhabricatorWorkerActiveTask`). Priority levels: ALERTS(1000) → DEFAULT(2000) → COMMIT(2500) → BULK(3000) → INDEX(3500) → IMPORT(4000). Tasks archived on completion. Yield/retry via exceptions (`PhabricatorWorkerYieldException`, `PhabricatorWorkerPermanentFailureException`).

**Harbormaster builds**: `HarbormasterBuildEngine` orchestrates build plans → `HarbormasterBuildWorker` → `HarbormasterTargetWorker` → step implementations (arc lint/unit, HTTP requests, Buildkite, CircleCI, Drydock commands).

**Drydock**: lease-based resource allocation for build resources. Allocator tries: free resources → new allocation → used resources → reclaim idle → yield and retry.

### Remarkup

Phabricator's markup language (Markdown-like with extensions). Batched rendering via `PhabricatorMarkupEngine` with cache-aware pipeline (preprocessed data cached in `PhabricatorMarkupCache`). Named engine variants per context (Maniphest, Phriction, Phame, Feed, Differential, Diffusion).

### Edge System

Typed directed PHID-to-PHID relationships stored per source object's database. `PhabricatorEdgeEditor` handles inverse edges automatically with cycle detection. `PhabricatorEdgeQuery` for reading (NOT policy-filtered — policy enforced when loading destination objects). Queries support edge-based constraints via `withEdgeLogicConstraints()`.

### Ferret Search

MySQL-based full-text search. Per-object-type schema: fdocument, ffield, fngrams, fngrams_common tables. Trigram indexing for substring matching. Function-scoped search (`title:foo`, `body:bar`). Integrated with query system via `withFerretConstraint()`.

## Naming Conventions

| Type | Pattern | Example |
|------|---------|---------|
| Application | `Phabricator{App}Application` | `PhabricatorManiphestApplication` |
| Controller | `{App}{Action}Controller` | `ManiphestTaskDetailController` |
| Storage | `{App}{Object}` | `ManiphestTask` |
| Query | `{App}{Object}Query` | `ManiphestTaskQuery` |
| Editor | `{App}{Object}Editor` | `ManiphestTransactionEditor` |
| Transaction type | `{App}{Object}{Field}Transaction` | `ManiphestTaskStatusTransaction` |
| Conduit | `{App}{Action}ConduitAPIMethod` | `ManiphestEditConduitAPIMethod` |

## Database Migrations

Add patches to `resources/sql/autopatches/` named `YYYYMMDD.{description}.{seq}.{sql|php}`. For declarative schema changes, update the app's `SchemaSpec` class and column types in `getConfiguration()`. Run `bin/storage upgrade` to apply.

## Library System

Phabricator is a Phutil library. The class map lives in `src/__phutil_library_map__.php` (12,800 lines) — regenerated by `arc liberate src/`. Must be run whenever classes are added, moved, or removed.

## PHP Compatibility

Minimum PHP 7.1. Active PHP 8 compatibility work in progress — recent commits address deprecations and breaking changes across Conduit, Dashboard, and miscellaneous applications.

## Coding Style

- 2-space indentation, 80-char line width
- No trailing whitespace, LF line endings
- PHP linted via xhpast (targeting PHP 5.5 syntax level) and phutil-library linters
- JS linted via jshint (browser config for webroot, node config for Aphlict)
- JS/CSS files require `@provides` docblock header for Celerity dependency tracking
