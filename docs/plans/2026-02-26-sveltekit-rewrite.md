# SvelteKit Rewrite Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rewrite Ghabricator from raw Go HTML rendering to a SvelteKit SPA with Go JSON API backend. Full feature parity.

**Architecture:** SvelteKit SPA (client-side only, adapter-static) consuming Go JSON API. Phabricator's PHUI CSS imported as-is. Javelin replaced by Svelte reactivity. Single Go binary embeds built frontend for production.

**Tech Stack:** SvelteKit 2, Svelte 5, TypeScript, Vite, bun, Go, Phabricator CSS, FontAwesome 4

---

## Phase 1: Foundation

### Task 1: Scaffold SvelteKit Project

**Files:**
- Create: `frontend/` (SvelteKit project)
- Create: `frontend/vite.config.ts`
- Create: `frontend/svelte.config.js`
- Create: `frontend/src/routes/+layout.ts`
- Create: `frontend/src/app.html`
- Create: `frontend/src/app.css`

**Step 1: Create SvelteKit project**

```bash
cd /Users/nikhilr/Code/experimental/phabricator/ghabricator
bunx sv create frontend --types ts --no-add-ons --install bun
```

Select: Skeleton project, TypeScript

**Step 2: Install adapter-static**

```bash
cd frontend && bun add -D @sveltejs/adapter-static
```

**Step 3: Configure svelte.config.js**

```js
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      fallback: 'index.html',
      precompress: false,
      strict: false,
    }),
  },
};

export default config;
```

**Step 4: Configure vite.config.ts with API proxy**

```ts
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

**Step 5: Disable SSR globally**

`frontend/src/routes/+layout.ts`:
```ts
export const ssr = false;
export const prerender = false;
```

**Step 6: Copy Phabricator CSS + FA4 into static/**

Copy from current `webroot/rsrc/` into `frontend/static/rsrc/`:
- All CSS files (phui/, application/, diff/, etc.)
- FontAwesome 4.7 (css + fonts)
- Phabricator fonts

Set up `app.html` to link the CSS bundle and FA4.

**Step 7: Verify dev server starts**

```bash
cd frontend && bun --bun run dev
```

Expected: SvelteKit skeleton running on :5173

**Step 8: Commit**

```bash
git add frontend/
git commit -m "feat: scaffold SvelteKit project with adapter-static and Phab CSS"
```

---

### Task 2: Go API Server Skeleton

**Files:**
- Modify: `internal/server/server.go` — add JSON API routes
- Create: `internal/server/api.go` — JSON response helpers
- Modify: `internal/server/server.go` — SPA fallback handler for production

**Step 1: Create JSON response helpers**

`internal/server/api.go`:
```go
package server

import (
    "encoding/json"
    "net/http"
)

func jsonOK(w http.ResponseWriter, data any) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
```

**Step 2: Add /api/auth/me endpoint**

Returns current user session as JSON. If not authenticated, returns 401.

```go
func (s *Server) handleAuthMe(w http.ResponseWriter, r *http.Request) {
    sess := s.auth.Store().GetFromRequest(r)
    if sess == nil {
        jsonError(w, "not authenticated", http.StatusUnauthorized)
        return
    }
    jsonOK(w, map[string]string{
        "login":     sess.Login,
        "avatarURL": sess.AvatarURL,
    })
}
```

**Step 3: Register API routes**

In `routes()`, add:
```go
// JSON API routes (SvelteKit frontend)
s.mux.Handle("GET /api/auth/me", s.auth.RequireAuth(http.HandlerFunc(s.handleAuthMe)))
```

Keep existing HTML routes for now — they'll coexist during migration.

**Step 4: Verify**

```bash
curl -s http://localhost:8080/api/auth/me | jq .
```

Expected: `{"error":"not authenticated"}` (401)

**Step 5: Commit**

```bash
git commit -m "feat: add Go JSON API skeleton with auth/me endpoint"
```

---

### Task 3: SvelteKit Auth Flow + Layout Shell

**Files:**
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/lib/stores/auth.ts`
- Create: `frontend/src/lib/stores/theme.ts`
- Create: `frontend/src/routes/+layout.svelte`
- Create: `frontend/src/lib/components/layout/Navbar.svelte`
- Create: `frontend/src/routes/+page.svelte` (landing)

**Step 1: Create API fetch wrapper**

`frontend/src/lib/api.ts`:
```ts
const BASE = '';  // same origin, proxied in dev

export async function apiFetch<T>(path: string, opts?: RequestInit): Promise<T> {
    const res = await fetch(`${BASE}${path}`, {
        credentials: 'include',
        ...opts,
    });
    if (res.status === 401) {
        window.location.href = '/api/auth/github';
        throw new Error('Unauthorized');
    }
    if (!res.ok) {
        const body = await res.json().catch(() => ({ error: res.statusText }));
        throw new Error(body.error || res.statusText);
    }
    return res.json();
}
```

**Step 2: Create auth store**

```ts
import { writable } from 'svelte/store';
import { apiFetch } from '$lib/api';

interface AuthUser {
    login: string;
    avatarURL: string;
}

export const user = writable<AuthUser | null>(null);
export const loading = writable(true);

export async function checkAuth() {
    try {
        const data = await apiFetch<AuthUser>('/api/auth/me');
        user.set(data);
    } catch {
        user.set(null);
    } finally {
        loading.set(false);
    }
}
```

**Step 3: Create theme store**

```ts
import { writable } from 'svelte/store';

function getInitialTheme(): string {
    if (typeof document === 'undefined') return '';
    const cookie = document.cookie.split(';').find(c => c.trim().startsWith('theme='));
    return cookie?.split('=')[1] === 'dark' ? 'dark' : '';
}

export const theme = writable(getInitialTheme());

export function toggleTheme() {
    theme.update(t => {
        const next = t === 'dark' ? '' : 'dark';
        document.cookie = `theme=${next};path=/;max-age=31536000;SameSite=Lax`;
        return next;
    });
}
```

**Step 4: Create Navbar component**

Renders the dark header bar with logo, nav links (Revisions, Repos, Paste, Herald, Search), user avatar, theme toggle. Uses same PHUI CSS classes as current Go template.

**Step 5: Create +layout.svelte**

```svelte
<script>
  import '../app.css';
  import Navbar from '$lib/components/layout/Navbar.svelte';
  import { checkAuth } from '$lib/stores/auth';
  import { onMount } from 'svelte';

  onMount(() => checkAuth());
</script>

<Navbar />
<slot />
```

**Step 6: Create landing page**

Same welcome content as current `handleIndex`, but in Svelte. If authenticated, redirect to `/dashboard`.

**Step 7: Verify**

```bash
cd frontend && bun --bun run dev
```

Open :5173 — should show landing page with navbar.

**Step 8: Commit**

```bash
git commit -m "feat: SvelteKit layout shell with auth, theme, and navbar"
```

---

## Phase 2: Go API Endpoints

Each task converts an existing HTML handler to return JSON instead.

### Task 4: Dashboard API

**Files:**
- Modify: `internal/server/dashboard.go` — add `handleAPIDashboard`

**Step 1: Create JSON dashboard handler**

Reuse existing `fetchAuthoredPRs` / `fetchReviewRequestedPRs` logic but return JSON:

```go
type DashboardResponse struct {
    Authored        []PRSummary `json:"authored"`
    ReviewRequested []PRSummary `json:"reviewRequested"`
}

type PRSummary struct {
    Number    int       `json:"number"`
    Title     string    `json:"title"`
    Owner     string    `json:"owner"`
    Repo      string    `json:"repo"`
    Author    User      `json:"author"`
    Draft     bool      `json:"draft"`
    Labels    []Label   `json:"labels"`
    Reviewers []User    `json:"reviewers"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

**Step 2: Register route**

```go
s.mux.Handle("GET /api/dashboard", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIDashboard)))
```

**Step 3: Verify**

```bash
curl -s http://localhost:8080/api/dashboard -b session_cookie | jq .
```

**Step 4: Commit**

---

### Task 5: PR Detail API

**Files:**
- Modify: `internal/server/pr_handler.go` — add `handleAPIPR`
- Modify: `internal/diff/parser.go` — ensure Changeset/Hunk/Line are JSON-serializable

The big one. Single endpoint returns everything the PR page needs.

**Step 1: Add JSON tags to diff types**

`internal/diff/types.go` — add `json:"..."` tags to LineType, Line, Hunk, Changeset, DiffRow.

**Step 2: Create PR detail response struct**

```go
type PRDetailResponse struct {
    PR            *ghapi.PullRequest    `json:"pr"`
    Changesets    []diff.Changeset      `json:"changesets"`
    Comments      map[string][]ghapi.ReviewComment `json:"commentsByPath"`
    Reviews       []ghapi.Review        `json:"reviews"`
    IssueComments []ghapi.IssueComment  `json:"issueComments"`
    CheckRuns     []ghapi.CheckRun      `json:"checkRuns"`
    Timeline      []TimelineEvent       `json:"timeline"`
    HeraldMatches []herald.RuleMatch    `json:"heraldMatches"`
}

type TimelineEvent struct {
    Author    ghapi.User `json:"author"`
    Action    string     `json:"action"`
    Body      string     `json:"body"`
    CreatedAt time.Time  `json:"createdAt"`
    IconClass string     `json:"iconClass"`
    IconColor string     `json:"iconColor"`
}
```

**Step 3: Create handler**

Same parallel fetch logic as existing `handlePR`, but returns JSON instead of rendering HTML. Syntax highlighting included in response as HTML strings per line (chroma output).

**Step 4: Add highlighted lines to response**

For each changeset, include pre-highlighted HTML per line so Svelte can render syntax-colored code without needing chroma on the client:

```go
type HighlightedChangeset struct {
    diff.Changeset
    Rows []diff.DiffRow `json:"rows"`
}
```

Call `diff.BuildDiffRows(cs)` + `diff.HighlightRows(cs)` to produce the rows with HTML content.

**Step 5: Register route and verify**

**Step 6: Commit**

---

### Task 6: Inline Comment API (already partially JSON)

**Files:**
- Modify: `internal/server/inline_handler.go`

Current inline handler already returns JSON (Javelin format). Simplify to standard JSON responses. Keep the same operations: new, save, edit, cancel, delete, done.

**Step 1: Refactor to return plain JSON**

Remove `for(;;);` Javelin prefix. Return standard `{"ok": true, "comment": {...}}` responses.

**Step 2: Commit**

---

### Task 7: Review / Merge / Close APIs

**Files:**
- Modify: `internal/server/pr_handler.go` (review, merge, close handlers)

These already accept POST. Change to return JSON instead of redirects.

**Step 1: Review handler — return JSON**

```go
jsonOK(w, map[string]string{"redirect": fmt.Sprintf("/pr/%s/%s/%d", owner, repo, number)})
```

**Step 2: Merge/close handlers — return JSON**

Same pattern.

**Step 3: Commit**

---

### Task 8: Repos API

**Files:**
- Modify: `internal/server/repos_handler.go` — add `handleAPIRepoList`
- Modify: `internal/server/repo_handler.go` — add `handleAPIRepoView`

**Step 1: Repo list API**

Returns `[]RepoSummary` with name, description, language, stars, forks, private, fork, archived, avatar, updatedAt.

**Step 2: Repo tree/file API**

Two endpoints or one smart endpoint:
- `GET /api/repo/:owner/:repo/tree?ref=X&path=Y` → directory entries
- `GET /api/repo/:owner/:repo/file?ref=X&path=Y` → file content (with syntax-highlighted lines for code, raw URL for images)
- `GET /api/repo/:owner/:repo/info` → repo metadata (for curtain)

**Step 3: Register routes and verify**

**Step 4: Commit**

---

### Task 9: Paste API

**Files:**
- Modify: `internal/server/paste_handler.go`

**Step 1: List pastes API** — `GET /api/paste` → `[]PasteSummary`

**Step 2: View paste API** — `GET /api/paste/:id` → `PasteDetail` with highlighted file contents

**Step 3: Create paste API** — `POST /api/paste` accepting JSON body → `{id, url}`

**Step 4: Register routes, verify, commit**

---

### Task 10: Herald API

**Files:**
- Modify: `internal/server/herald_handler.go`

Herald types already have JSON tags. Mostly just register CRUD endpoints:

- `GET /api/herald` → `[]Rule`
- `GET /api/herald/:id` → `Rule`
- `POST /api/herald` → `Rule` (create/update)
- `DELETE /api/herald/:id` → `{ok: true}`

**Step 1: Implement handlers**
**Step 2: Register routes, verify, commit**

---

### Task 11: Search API

**Files:**
- Modify: `internal/server/search_handler.go`

`GET /api/search?q=X&type=Y` → returns typed results (PRs, code, repos).

**Step 1: Implement handler returning JSON**
**Step 2: Register route, verify, commit**

---

## Phase 3: Svelte PHUI Component Library

These components wrap Phabricator CSS classes in Svelte components with typed props.

### Task 12: Layout Components

**Files:**
- Create: `frontend/src/lib/components/layout/PageShell.svelte`
- Create: `frontend/src/lib/components/layout/Breadcrumbs.svelte`
- Create: `frontend/src/lib/components/layout/CurtainLayout.svelte`
- Create: `frontend/src/lib/components/layout/FormationView.svelte`

**PageShell** — main content area with optional header (icon + title + subtitle), wraps `<slot>`.

**Breadcrumbs** — accepts `crumbs: {name: string, href?: string}[]`, renders chevron-separated links.

**CurtainLayout** — two-column layout: `<slot>` (main) + `<slot name="curtain">` (sidebar).

**FormationView** — three-column layout with collapsible left sidebar (file tree), main content, right curtain.

All using existing PHUI CSS classes from the Phabricator stylesheets.

**Commit after each component or batch.**

---

### Task 13: PHUI Components

**Files:**
- Create: `frontend/src/lib/components/phui/ObjectItem.svelte`
- Create: `frontend/src/lib/components/phui/ObjectItemList.svelte`
- Create: `frontend/src/lib/components/phui/Tag.svelte`
- Create: `frontend/src/lib/components/phui/Button.svelte`
- Create: `frontend/src/lib/components/phui/InfoView.svelte`
- Create: `frontend/src/lib/components/phui/HeaderView.svelte`
- Create: `frontend/src/lib/components/phui/PropertyList.svelte`
- Create: `frontend/src/lib/components/phui/Timeline.svelte`
- Create: `frontend/src/lib/components/phui/StatusList.svelte`
- Create: `frontend/src/lib/components/phui/ActionList.svelte`
- Create: `frontend/src/lib/components/phui/Box.svelte`

Each component:
1. Accepts typed props matching the PHUI CSS vocabulary
2. Renders the correct `class="phui-*"` markup
3. Exposes slots for flexible content

Example — `Tag.svelte`:
```svelte
<script lang="ts">
  export let shade: string = 'blue';
  export let icon: string = '';
</script>

<span class="phui-tag-view phui-tag-shade-{shade} phui-tag-type-shade">
  <span class="phui-tag-core">
    {#if icon}<span class="phui-icon-view phui-font-fa {icon} mrs"></span>{/if}
    <slot />
  </span>
</span>
```

**Commit after batch.**

---

### Task 14: Diff Components

The most complex component set. Replaces Go's `RenderChangeset` + Javelin behaviors.

**Files:**
- Create: `frontend/src/lib/components/diff/DiffTable.svelte`
- Create: `frontend/src/lib/components/diff/ChangesetHeader.svelte`
- Create: `frontend/src/lib/components/diff/FileTree.svelte`
- Create: `frontend/src/lib/components/diff/InlineComment.svelte`
- Create: `frontend/src/lib/components/diff/InlineEditor.svelte`
- Create: `frontend/src/lib/components/diff/ContextExpander.svelte`
- Create: `frontend/src/lib/stores/inline.ts`

**DiffTable** — renders side-by-side diff from `DiffRow[]` data. Each row has old/new line numbers, old/new content (pre-highlighted HTML from API), and CSS classes. Renders as `<table class="phabricator-source-code-view">`.

**ChangesetHeader** — sticky file header with icon, path, +/- stats, collapse toggle.

**FileTree** — left sidebar with file list, click-to-scroll, view progress tracking.

**InlineComment** — displays comment with avatar, author, body, actions (Reply, Done).

**InlineEditor** — textarea for new/edit inline comments. Saves via `POST /api/inline`.

**ContextExpander** — "show more lines" row. Clicks fetch `/api/context` and splice results into the rows reactively.

**inline.ts store** — tracks pending inline comments, drafts, done states. Replaces Javelin's `differential-populate` behavior.

**Commit after each component or logical batch.**

---

### Task 15: Review Form Component

**Files:**
- Create: `frontend/src/lib/components/review/ReviewForm.svelte`

Textarea + 3 buttons (Comment, Accept, Request Changes). Shows pending inline comment count from store. Posts to `/api/review`.

**Commit.**

---

## Phase 4: Svelte Pages

Each page composes the components from Phase 3 with data from Phase 2 APIs.

### Task 16: Dashboard Page

**Files:**
- Create: `frontend/src/routes/dashboard/+page.svelte`
- Create: `frontend/src/routes/dashboard/+page.ts`

**+page.ts** loads data:
```ts
import { apiFetch } from '$lib/api';
export async function load() {
    return { dashboard: await apiFetch('/api/dashboard') };
}
```

**+page.svelte** renders:
- PageShell with "Active Revisions" header
- "Authored" section: ObjectItemList with PR items (avatar, title, tags, attributes)
- "Review Requested" section: same format
- Empty states with InfoView

**Commit.**

---

### Task 17: PR Detail Page

**Files:**
- Create: `frontend/src/routes/pr/[owner]/[repo]/[number]/+page.svelte`
- Create: `frontend/src/routes/pr/[owner]/[repo]/[number]/+page.ts`

The biggest page. Uses FormationView (3-column):
- Left: FileTree
- Center: Summary box + Changesets (DiffTable per file) + Timeline + ReviewForm
- Right: Curtain panels (Reviewers, Buildables, Herald, Properties, Actions)

**+page.ts** loads single `/api/pr/:owner/:repo/:number`.

**+page.svelte** composes all diff components, timeline, review form, curtain panels.

**Commit.**

---

### Task 18: Repos Pages

**Files:**
- Create: `frontend/src/routes/repos/+page.svelte` (list)
- Create: `frontend/src/routes/repos/[owner]/[repo]/[...path]/+page.svelte` (browser)

**List:** ObjectItemList with repo items (avatar, name, tags, language dot, stars, forks, time).

**Browser:** Smart page that fetches tree or file based on path. Dir view = table. File view = syntax-highlighted code or image preview. Curtain with repo info.

**Commit.**

---

### Task 19: Paste Pages

**Files:**
- Create: `frontend/src/routes/paste/+page.svelte` (list)
- Create: `frontend/src/routes/paste/new/+page.svelte` (create form)
- Create: `frontend/src/routes/paste/[id]/+page.svelte` (view)

**Commit.**

---

### Task 20: Herald Pages

**Files:**
- Create: `frontend/src/routes/herald/+page.svelte` (list)
- Create: `frontend/src/routes/herald/new/+page.svelte` (create form)
- Create: `frontend/src/routes/herald/[id]/+page.svelte` (view/edit)

Dynamic condition/action row management is much cleaner in Svelte than the current inline JS.

**Commit.**

---

### Task 21: Search Page

**Files:**
- Create: `frontend/src/routes/search/+page.svelte`

Search form + typed results sections (PRs, code, repos). Reuses ObjectItemList.

**Commit.**

---

## Phase 5: Integration & Cleanup

### Task 22: Dark Mode

**Files:**
- Modify: `frontend/src/routes/+layout.svelte` — apply theme class to body
- Modify: `frontend/src/app.css` — import dark mode CSS overrides

Body gets `phui-theme-blindigo` always, plus `phui-theme-dark` when dark. Same CSS scoping as current implementation.

**Commit.**

---

### Task 23: Go Embed + Production Build

**Files:**
- Create: `frontend_embed.go` — `//go:embed frontend/build`
- Modify: `internal/server/server.go` — SPA fallback handler
- Modify: `cmd/ghabricator/main.go` — wire embed

**Step 1: Add embed directive**

```go
//go:embed frontend/build
var FrontendFS embed.FS
```

**Step 2: Add SPA handler**

Try to serve static file from embedded FS; fall back to `index.html` for SvelteKit client routes.

**Step 3: Wire routing**

API routes match first (longer prefix), everything else falls through to SPA handler.

**Step 4: Build and verify**

```bash
cd frontend && bun run build && cd ..
go build -o /tmp/ghabricator ./cmd/ghabricator
/tmp/ghabricator
```

Single binary serves both API and frontend on :8080.

**Step 5: Commit**

---

### Task 24: Remove Old HTML Handlers

**Files:**
- Remove HTML rendering from all handlers (keep only JSON API versions)
- Remove: `internal/templates/page.html`
- Remove: `internal/templates/render.go`
- Remove: old `webroot/rsrc/` (now in `frontend/static/rsrc/`)
- Clean up imports

**Commit.**

---

## Execution Order & Parallelism

```
Phase 1 (sequential):  Task 1 → Task 2 → Task 3

Phase 2 + 3 (parallel):
  Track A (Go API):     Tasks 4-11
  Track B (Components): Tasks 12-15

Phase 4 (sequential, depends on 2+3): Tasks 16-21

Phase 5 (sequential): Tasks 22-24
```

**Estimated tasks: 24**
**Parallelizable: Phase 2 and Phase 3 can run simultaneously**
