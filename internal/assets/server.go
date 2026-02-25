package assets

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikhilr/ghabricator/internal/diff"
)

// Server serves Phabricator's static assets with CSS variable replacement.
type Server struct {
	// webroot is the path to the webroot/rsrc/ directory
	webroot string
	// Built packages: "core.pkg.css" → processed bytes
	packages map[string][]byte
	// ETags for packages
	etags map[string]string
}

// NewServer builds the asset server, parsing the celerity map and pre-building packages.
func NewServer(repoRoot string) (*Server, error) {
	webroot := filepath.Join(repoRoot, "webroot", "rsrc")
	mapPath := filepath.Join(repoRoot, "resources", "celerity", "map.php")

	cm, err := ParseCelerityMap(mapPath)
	if err != nil {
		return nil, err
	}

	s := &Server{
		webroot:  webroot,
		packages: make(map[string][]byte),
		etags:    make(map[string]string),
	}

	// Build packages we care about
	pkgNames := []string{
		"core.pkg.css",
		"core.pkg.js",
		"differential.pkg.css",
		"differential.pkg.js",
	}

	for _, pkg := range pkgNames {
		paths, err := cm.ResolvePackage(pkg)
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}

		isCSS := strings.HasSuffix(pkg, ".css")
		var buf []byte

		for _, p := range paths {
			// Paths in map.php are relative like "rsrc/css/core/core.css"
			fullPath := filepath.Join(repoRoot, "webroot", p)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("warning: package %s: skip %s: %v", pkg, p, err)
				continue
			}
			if isCSS {
				data = ProcessCSS(data, DefaultTheme)
			}
			buf = append(buf, data...)
			buf = append(buf, '\n')
		}

		// Append Chroma syntax highlighting CSS to CSS packages
		if isCSS {
			buf = append(buf, []byte("\n"+diff.ChromaCSSLight+"\n"+diff.ChromaCSSDark)...)
		}

		s.packages[pkg] = buf
		s.etags[pkg] = fmt.Sprintf(`"%x"`, sha256.Sum256(buf))

		// Also build dark mode variant for CSS packages
		if isCSS {
			darkKey := "dark/" + pkg
			var darkBuf []byte
			for _, p := range paths {
				fullPath := filepath.Join(repoRoot, "webroot", p)
				data, err := os.ReadFile(fullPath)
				if err != nil {
					continue
				}
				data = ProcessCSS(data, DarkTheme)
				darkBuf = append(darkBuf, data...)
				darkBuf = append(darkBuf, '\n')
			}
			darkBuf = append(darkBuf, []byte("\n"+diff.ChromaCSSLight+"\n"+diff.ChromaCSSDark)...)
			s.packages[darkKey] = darkBuf
			s.etags[darkKey] = fmt.Sprintf(`"%x"`, sha256.Sum256(darkBuf))
		}

		log.Printf("built package %s (%d bytes, %d files)", pkg, len(buf), len(paths))
	}

	return s, nil
}

// ServeHTTP handles GET /res/ requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip /res/ prefix
	path := strings.TrimPrefix(r.URL.Path, "/res/")

	// Check for package requests: /res/pkg/core.pkg.css or /res/pkg/dark/core.pkg.css
	if strings.HasPrefix(path, "pkg/") {
		pkgKey := strings.TrimPrefix(path, "pkg/")
		if data, ok := s.packages[pkgKey]; ok {
			etag := s.etags[pkgKey]
			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			w.Header().Set("ETag", etag)
			if strings.HasSuffix(pkgKey, ".css") {
				w.Header().Set("Content-Type", "text/css; charset=utf-8")
			} else {
				w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			}
			w.Write(data)
			return
		}
		http.NotFound(w, r)
		return
	}

	// Serve individual files from webroot/rsrc/
	fullPath := filepath.Join(s.webroot, filepath.Clean(path))

	// Security: don't allow path traversal outside webroot
	if !strings.HasPrefix(fullPath, s.webroot) {
		http.NotFound(w, r)
		return
	}

	// For CSS files, apply variable replacement
	if strings.HasSuffix(path, ".css") {
		data, err := os.ReadFile(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		data = ProcessCSS(data, DefaultTheme)
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Write(data)
		return
	}

	// Everything else: serve as-is
	http.ServeFile(w, r, fullPath)
}
