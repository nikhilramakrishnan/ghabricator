# SvelteKit Rewrite Design

## Summary

Rewrite Ghabricator's frontend from raw Go HTML string concatenation to a SvelteKit SPA, keeping Go as a JSON API backend. Full feature parity with all current pages. Phabricator's PHUI CSS imported as-is; Javelin replaced by Svelte reactivity.

## Architecture

```
┌─────────────────────────────┐
│  SvelteKit SPA (frontend/)  │
│  Routes, components, stores │
│  Phabricator CSS imported   │
├─────────────────────────────┤
│         /api/* proxy        │
├─────────────────────────────┤
│  Go API server (cmd/api)    │
│  JSON endpoints, OAuth,     │
│  Herald, diff parsing       │
│  Embeds built SvelteKit     │
└─────────────────────────────┘
```

- **Dev:** SvelteKit on :5173, proxies `/api/*` to Go on :8080
- **Prod:** Go serves built SvelteKit from `go:embed` + API on same port

## Go API Surface

All current handler logic stays, returns JSON instead of HTML.

| Endpoint | Method | Returns |
|---|---|---|
| `/api/auth/github` | GET | Redirect (OAuth initiate) |
| `/api/auth/callback` | GET | Sets session cookie, redirects |
| `/api/auth/logout` | GET | Clears session |
| `/api/auth/me` | GET | Current user (login, avatar) |
| `/api/dashboard` | GET | Authored PRs + review-requested PRs |
| `/api/pr/:owner/:repo/:number` | GET | PR detail + diff + comments + reviews + checks |
| `/api/review` | POST | Submit review |
| `/api/inline` | POST | Inline comment CRUD |
| `/api/merge` | POST | Merge PR |
| `/api/close` | POST | Close/reopen PR |
| `/api/context` | GET | Diff context expansion |
| `/api/repos` | GET | User's repositories |
| `/api/repo/:owner/:repo/tree` | GET | Directory listing |
| `/api/repo/:owner/:repo/file` | GET | File content |
| `/api/paste` | GET | List pastes |
| `/api/paste` | POST | Create paste |
| `/api/paste/:id` | GET | View paste |
| `/api/herald` | GET | List rules |
| `/api/herald` | POST | Create/update rule |
| `/api/herald/:id` | GET | View rule |
| `/api/herald/:id` | DELETE | Delete rule |
| `/api/search` | GET | Search (PRs, code, repos) |

The PR detail endpoint does parallel fetches internally (PR metadata, diff, comments, reviews, checks) and returns a single JSON blob.

## SvelteKit Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api.ts              # Fetch wrapper (handles auth, errors)
│   │   ├── stores/
│   │   │   ├── auth.ts         # User session store
│   │   │   └── theme.ts        # Dark/light mode
│   │   └── components/
│   │       ├── layout/
│   │       │   ├── Navbar.svelte
│   │       │   ├── Breadcrumbs.svelte
│   │       │   └── CurtainPanel.svelte
│   │       ├── phui/            # Phabricator component wrappers
│   │       │   ├── ObjectItem.svelte
│   │       │   ├── ObjectItemList.svelte
│   │       │   ├── Tag.svelte
│   │       │   ├── Button.svelte
│   │       │   ├── InfoView.svelte
│   │       │   ├── HeaderView.svelte
│   │       │   ├── PropertyList.svelte
│   │       │   ├── Timeline.svelte
│   │       │   ├── StatusList.svelte
│   │       │   └── ActionList.svelte
│   │       ├── diff/
│   │       │   ├── DiffTable.svelte
│   │       │   ├── ChangesetHeader.svelte
│   │       │   ├── FileTree.svelte
│   │       │   ├── InlineComment.svelte
│   │       │   └── InlineEditor.svelte
│   │       └── review/
│   │           └── ReviewForm.svelte
│   ├── routes/
│   │   ├── +layout.svelte      # Shell: navbar, theme, auth check
│   │   ├── +page.svelte        # Landing/login
│   │   ├── dashboard/
│   │   │   └── +page.svelte
│   │   ├── pr/[owner]/[repo]/[number]/
│   │   │   └── +page.svelte
│   │   ├── repos/
│   │   │   ├── +page.svelte
│   │   │   └── [owner]/[repo]/[...path]/
│   │   │       └── +page.svelte
│   │   ├── paste/
│   │   │   ├── +page.svelte
│   │   │   ├── new/+page.svelte
│   │   │   └── [id]/+page.svelte
│   │   ├── herald/
│   │   │   ├── +page.svelte
│   │   │   ├── new/+page.svelte
│   │   │   └── [id]/+page.svelte
│   │   └── search/
│   │       └── +page.svelte
│   └── app.css                 # Phabricator CSS imports + overrides
├── static/
│   └── rsrc/                   # Phab CSS/fonts copied here
├── svelte.config.js
├── vite.config.ts
└── package.json
```

## PHUI Component Library

Svelte wrappers around Phabricator's PHUI markup. Instead of raw HTML string concatenation, components accept props and render the correct PHUI classes:

```svelte
<ObjectItem title={pr.title} href={prUrl} barColor="blue">
  <Tag slot="tags" shade="green" icon="fa-check">Accepted</Tag>
  <svelte:fragment slot="attributes">
    <Attribute icon="fa-user">{pr.author}</Attribute>
    <Attribute icon="fa-clock-o">{timeAgo(pr.updatedAt)}</Attribute>
  </svelte:fragment>
</ObjectItem>
```

One line to add a tag or attribute, not 5 lines of escaped HTML. This is where the richness comes from.

## Data Flow

**Auth:** SvelteKit `+layout.ts` checks `/api/auth/me` on load. 401 = redirect to login. Session cookie managed by Go.

**PR detail:** Single API call returns everything. Diff rendering in `DiffTable.svelte`. Inline comments tracked in Svelte stores with optimistic updates. Context expansion via reactive splice. Review form with live pending comment count.

**Theme:** Svelte store + cookie. `<body class:phui-theme-dark={$theme === 'dark'}>`.

**Herald:** Pure CRUD forms. Go stores rules in memory.

## What Changes, What Stays

| Layer | Now | After |
|---|---|---|
| Routing | Go `http.ServeMux` | SvelteKit file routes |
| Page rendering | Go `WriteString` raw HTML | Svelte components |
| Interactivity | Javelin behaviors | Svelte reactivity |
| Diff rendering | Go HTML output | Svelte `DiffTable` component |
| Diff parsing | Go `diff.ParseDiff` | **Stays in Go**, returns JSON |
| Syntax highlighting | Go chroma | **Stays in Go**, returns HTML spans |
| Auth | Go OAuth + cookies | **Stays in Go**, SvelteKit reads cookie |
| Herald | Go in-memory store | **Stays in Go**, JSON API |
| CSS | Phab CSS served by Go | Phab CSS imported into SvelteKit |
| Dark mode | Cookie + Go template class | Cookie + Svelte store |

## Pages (Full Parity)

1. Landing/login
2. Dashboard (authored + review-requested PRs)
3. PR detail (diffs, inline comments, timeline, review form, curtain with reviewers/buildables/herald/properties/actions)
4. Repository list
5. Repository browser (directory + file views with syntax highlighting + image preview)
6. Paste list / create / view
7. Herald list / create / view
8. Search (PRs, code, repos)
