# CLAUDE.md

## What This Is

Ghabricator — a modern GitHub code review UI inspired by Phabricator. Go backend proxies to GitHub's API; SvelteKit frontend renders a Phabricator-style experience (side-by-side diffs, inline comments, Herald automation, repository browsing).

No PHP, no database — GitHub is the source of truth.

## Environment & Dependencies

- **Go 1.22+** for the backend
- **bun** for the frontend (not npm)
- **Node.js v22** available
- No PHP, MySQL, or Arcanist needed

### Go Dependencies (`go.mod`)
- `github.com/google/go-github/v68` — GitHub API client
- `github.com/alecthomas/chroma/v2` — syntax highlighting
- `github.com/sourcegraph/go-diff` — diff parsing
- `github.com/yuin/goldmark` — markdown rendering
- `golang.org/x/oauth2` — OAuth2 client

### Frontend Dependencies
- SvelteKit 2.50, Svelte 5.51, TypeScript
- `adapter-static` (SPA mode, no SSR)

### Required Environment Variables
- `GITHUB_CLIENT_ID` — GitHub OAuth app client ID
- `GITHUB_CLIENT_SECRET` — GitHub OAuth app secret
- `SESSION_SECRET` — HMAC signing key (defaults to dev secret)
- `PORT` — HTTP port (default: 8080)

## Development Commands

```bash
# Backend
go build -o ghabricator ./cmd/ghabricator    # Build
PORT=8080 ./ghabricator                       # Run

# Frontend
cd frontend && bun install                    # Install deps
cd frontend && bun run dev                    # Dev server on :5173
cd frontend && bun run build                  # Production build → build/
```

## Source Layout

```
cmd/ghabricator/main.go          # Entry point — HTTP server
internal/
  server/                        # Route setup, all API handlers
    server.go                    # Server struct, routes, middleware
    dashboard.go                 # GET /api/dashboard
    pr.go                        # GET /api/pr/{o}/{r}/{n}
    inline.go                    # POST /api/v2/inline (drafts + GitHub comments)
    review.go                    # POST /api/v2/review, /merge, /close
    repos.go                     # GET /api/repos, /repo/{o}/{r}/info|tree|file
    paste.go                     # GET/POST /api/paste
    herald.go                    # CRUD /api/herald
    search.go                    # GET /api/search
  auth/                          # GitHub OAuth2 + session management
    oauth.go                     # OAuth flow handlers
    session.go                   # Cookie sessions (HMAC-signed, in-memory store)
  github/                        # GitHub API wrapper
    client.go                    # PR, repo, gist, search operations
    types.go                     # Go types for GitHub objects
  diff/                          # Diff parsing + syntax highlighting
    parse.go                     # ParseDiff → []Changeset with []DiffRow
    highlight.go                 # HighlightLines via chroma
  herald/                        # Herald automation rules
    engine.go                    # Evaluate(rules, PRContext) → []RuleMatch
    store.go                     # JSON file storage (~/.ghabricator/herald-rules.json)
    types.go                     # Rule, Condition, Action types
frontend/
  src/
    routes/                      # SvelteKit file-based routes
      +layout.svelte             # Auth check, theme, navbar
      dashboard/                 # Authored + review-requested PRs
      pr/[owner]/[repo]/[number] # PR detail with diff viewer
      repos/                     # Repository list + file browser
      search/                    # Search (PRs, code, repos)
      paste/                     # Gist list, view, create
      herald/                    # Herald rules CRUD
    lib/
      api.ts                     # apiFetch/apiPost helpers
      types.ts                   # TypeScript types mirroring Go structs
      stores/                    # auth, theme, inline draft stores
      components/                # PHUI components, diff viewer, review form
```

## API Surface

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/auth/me` | GET | Current user info |
| `/api/dashboard` | GET | Authored + review-requested PRs |
| `/api/pr/{o}/{r}/{n}` | GET | Full PR detail (parallel fetches) |
| `/api/v2/inline` | POST | Inline comment CRUD (new/save/edit/delete/cancel/done) |
| `/api/v2/review` | POST | Submit review (approve/request changes/comment) |
| `/api/v2/merge` | POST | Merge PR (merge/squash/rebase) |
| `/api/v2/close` | POST | Close/reopen PR |
| `/api/repos` | GET | List user repos |
| `/api/repo/{o}/{r}/info` | GET | Repo metadata |
| `/api/repo/{o}/{r}/tree` | GET | Directory listing |
| `/api/repo/{o}/{r}/file` | GET | File content + syntax highlighting |
| `/api/paste` | GET/POST | List or create gists |
| `/api/paste/{id}` | GET | View gist |
| `/api/herald` | GET/POST | List or save rules |
| `/api/herald/{id}` | GET/DELETE | Get or delete rule |
| `/api/search` | GET | Search PRs/code/repos |
| `/auth/github` | GET | OAuth login |
| `/auth/callback` | GET | OAuth callback |
| `/auth/logout` | GET | Logout |

## Architecture

### Auth Flow
`/auth/github` → GitHub OAuth authorize → `/auth/callback` → exchange code for token → create session → redirect to dashboard. Sessions stored in-memory with HMAC-signed cookies (`phab_session`, 24h TTL). `RequireAuth` middleware injects session + GitHub client into request context.

### Request Flow
```
HTTP request → ServeMux route match → RequireAuth middleware
  → Handler: extract session → build GitHub client from token
  → Fetch data from GitHub API (parallel where possible)
  → Transform to API response structs → JSON response
```

PR detail handler does 5 parallel GitHub API calls (PR metadata, diff, review comments, reviews, issue comments) then a sequential check runs fetch.

### Diff System
Raw unified diffs from GitHub → `sourcegraph/go-diff` parser → `[]Changeset` with side-by-side `[]DiffRow` → syntax highlighted via `chroma/v2`. Comments indexed by file path for O(1) lookup.

### Herald Rules
JSON file at `~/.ghabricator/herald-rules.json`. Conditions: file_path (glob), author, title, label, base_branch. Actions: add_reviewer, add_label, post_comment. Evaluated per-PR with AND/OR logic.

### Inline Comments
In-memory draft store (`map[int64]*inlineDraft` with mutex). Operations: `new` (create draft) → `save` (POST to GitHub) → `edit`/`delete`/`cancel`/`done`.

## Frontend Architecture

### Component Hierarchy
- **Layout:** PageShell, CurtainLayout (two-column), Navbar, Breadcrumbs
- **PHUI (Phabricator-style):** Box, CurtainBox, Button, ActionList, ObjectItem, PropertyList, StatusList, Tag, Timeline, HeaderView, InfoView
- **Diff:** DiffTable (side-by-side), ChangesetHeader, InlineComment, InlineEditor, ContextExpander, FileTree
- **Review:** ReviewForm (approve/request changes/comment)

### Data Layer
- `lib/api.ts` — `apiFetch<T>()` / `apiPost<T>()` with credentials and 401 redirect
- `lib/stores/auth.ts` — user session state
- `lib/stores/theme.ts` — dark/light mode (cookie-based)
- `lib/stores/inline.ts` — draft comment tracking

### CSS
All scoped Svelte styles. Dark mode via body class toggle. No preprocessor.

## Coding Style

- Go: standard `gofmt`, idiomatic patterns
- TypeScript/Svelte: 2-space indent, strict types
- Lean code, no overengineering
- Icons over text in UI
