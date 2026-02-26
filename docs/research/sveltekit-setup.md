# SvelteKit + Vite + Phab CSS Setup Research

## 1. Scaffolding with Bun

```bash
cd /Users/nikhilr/Code/experimental/phabricator/ghabricator
bun create svelte@latest frontend
```

Prompts to answer:
- Template: **Skeleton project** (minimal)
- Type checking: **TypeScript**
- Add-ons: none (skip)

Then:
```bash
cd frontend
bun install
bun --bun run dev   # starts on :5173
```

Alternatively, the newer `sv` CLI:
```bash
bunx sv create frontend --types ts --no-add-ons --install bun
```

Install adapter-static:
```bash
cd frontend
bun add -D @sveltejs/adapter-static
```

## 2. Vite Proxy Config

**vite.config.ts** — proxy `/api/*` to Go backend on :8080 during dev:

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

No `rewrite` needed — we want `/api/auth/me` to hit `http://localhost:8080/api/auth/me` as-is.

**Note:** This proxy only applies to the Vite dev server. In production, Go serves both API and static files on the same port.

## 3. Importing External CSS into SvelteKit

### Option A: Import in +layout.svelte (recommended — gets HMR)

```svelte
<!-- src/routes/+layout.svelte -->
<script>
  import '../app.css';
</script>

<slot />
```

**src/app.css:**
```css
/* Import Phabricator CSS files from static/rsrc/ */
@import '/rsrc/css/phui/phui-theme.css';
@import '/rsrc/css/phui/phui-button.css';
@import '/rsrc/css/phui/phui-object-item.css';
/* ... etc */

/* Or if copying CSS into src/ for Vite processing: */
@import './phab/phui-theme.css';
```

### Option B: Link in app.html (no HMR, but simpler for bulk CSS)

```html
<!-- src/app.html -->
<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="stylesheet" href="/rsrc/css/all-phab.css" />
  %sveltekit.head%
</head>
<body>
  %sveltekit.body%
</body>
</html>
```

### Option C: Place CSS in static/ directory

Files in `static/` are served as-is at the root URL. Put Phab CSS at `static/rsrc/css/` and reference via `<link>` or `@import url('/rsrc/css/...')`.

**Recommendation:** Use Option C for bulk Phab CSS (copy files into `static/rsrc/`), with an `app.css` that has any overrides or custom styles imported via Option A.

## 4. SvelteKit as Pure Client-Side SPA (No SSR)

### svelte.config.js

```js
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      fallback: 'index.html',  // SPA fallback — all routes serve this
      precompress: false,
      strict: false,
    }),
  },
};

export default config;
```

### src/routes/+layout.ts (or +layout.js)

```ts
// Disable SSR globally — pure client-side SPA
export const ssr = false;
export const prerender = false;
```

### Build

```bash
bun run build    # outputs to frontend/build/
```

Output structure:
```
build/
├── _app/           # Hashed JS/CSS chunks
│   ├── immutable/
│   │   ├── chunks/
│   │   ├── entry/
│   │   └── nodes/
│   └── version.json
├── index.html      # SPA fallback (all routes)
└── favicon.png
```

## 5. Go Embed + SPA Fallback Routing

### Directory structure

```
ghabricator/
├── cmd/api/main.go
├── frontend/          # SvelteKit project
│   ├── build/         # Built output (git-ignored)
│   └── ...
├── embed.go           # Embed directive
└── ...
```

### embed.go

```go
package ghabricator

import "embed"

//go:embed frontend/build
var FrontendFS embed.FS
```

### SPA handler in cmd/api/main.go

```go
package main

import (
    "embed"
    "io/fs"
    "net/http"
    "strings"
)

func spaHandler(embeddedFS embed.FS, buildDir string) http.Handler {
    fsys, err := fs.Sub(embeddedFS, buildDir)
    if err != nil {
        panic(err)
    }
    fileServer := http.FileServer(http.FS(fsys))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Try to serve the file directly
        path := strings.TrimPrefix(r.URL.Path, "/")
        if path == "" {
            path = "index.html"
        }

        // Check if file exists in embedded FS
        f, err := fsys.Open(path)
        if err != nil {
            // File not found — serve index.html (SPA fallback)
            r.URL.Path = "/"
            fileServer.ServeHTTP(w, r)
            return
        }
        f.Close()

        // File exists — serve it directly
        fileServer.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()

    // API routes
    mux.HandleFunc("/api/", apiHandler)

    // SPA fallback — all non-API routes serve SvelteKit
    mux.Handle("/", spaHandler(FrontendFS, "frontend/build"))

    http.ListenAndServe(":8080", mux)
}
```

**Key points:**
- API routes (`/api/*`) are matched first by Go's `ServeMux` (longer prefix wins)
- Everything else falls through to the SPA handler
- SPA handler tries to serve the exact file (JS, CSS, images) — if not found, serves `index.html`
- The `//go:embed frontend/build` directive bundles the entire build output into the binary
- For dev mode, skip the embed and just proxy — SvelteKit dev server handles everything

### Build command

```bash
cd frontend && bun run build && cd ..
go build -o ghabricator ./cmd/api
# Single binary, serves both API and frontend
```

## 6. FontAwesome 4

Phabricator uses FontAwesome 4.7.0 (`fa fa-check`, `fa fa-user`, etc.).

### Option A: CDN (simplest)

In `src/app.html`:
```html
<link rel="stylesheet"
  href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css"
  integrity="sha512-SfTiTlX6kk+qitfevl/7LibUOeJWlt9rbyDn92a1DqWOw9vwG2MFoQvAd8q62aA3J2xvgFlEHQ+LetZcIYR/Ew=="
  crossorigin="anonymous" />
```

### Option B: Self-hosted (for embed in Go binary)

Download FA 4.7.0 and place in `static/`:
```
static/
└── rsrc/
    └── font-awesome/
        ├── css/
        │   └── font-awesome.min.css
        └── fonts/
            ├── fontawesome-webfont.woff2
            ├── fontawesome-webfont.woff
            ├── fontawesome-webfont.ttf
            └── FontAwesome.otf
```

Fix font paths in `font-awesome.min.css` (change `../fonts/` to `/rsrc/font-awesome/fonts/`).

Then in `src/app.html`:
```html
<link rel="stylesheet" href="/rsrc/font-awesome/css/font-awesome.min.css" />
```

**Recommendation:** Option B (self-hosted) — the Go binary embeds everything, no CDN dependency. The FA4 CSS + fonts total ~200KB.

### Usage in Svelte

```svelte
<i class="fa fa-check"></i>
<i class="fa fa-user"></i>
<i class="fa fa-code"></i>
```

Same classes as Phabricator — no changes needed.

## Quick Reference: Complete Config Files

### package.json (key deps)
```json
{
  "devDependencies": {
    "@sveltejs/adapter-static": "^3.0.0",
    "@sveltejs/kit": "^2.0.0",
    "svelte": "^5.0.0",
    "vite": "^6.0.0"
  }
}
```

### tsconfig.json
Generated by `sv create`. No changes needed.

### Dev workflow
```bash
# Terminal 1: Go API
go run ./cmd/api -dev

# Terminal 2: SvelteKit dev server
cd frontend && bun --bun run dev
```

SvelteKit on :5173 proxies `/api/*` to Go on :8080. Browser hits :5173.

### Prod build
```bash
cd frontend && bun run build && cd ..
go build -o ghabricator ./cmd/api
./ghabricator  # serves everything on :8080
```
