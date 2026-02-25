package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
)

// handleContext serves expanded diff context lines via AJAX.
// GET /api/context?owner=X&repo=X&ref=X&path=X&start=N&end=N&cs=N
func (s *Server) handleContext(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())

	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	ref := r.URL.Query().Get("ref")
	path := r.URL.Query().Get("path")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	csStr := r.URL.Query().Get("cs")

	start, _ := strconv.Atoi(startStr)
	end, _ := strconv.Atoi(endStr)
	csID, _ := strconv.Atoi(csStr)

	if owner == "" || repo == "" || path == "" || start == 0 || end == 0 {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}
	if ref == "" {
		ref = "HEAD"
	}

	file, err := ghapi.FetchFileContent(r.Context(), client, owner, repo, ref, path)
	if err != nil {
		http.Error(w, "failed to fetch file", http.StatusInternalServerError)
		return
	}

	html := diff.RenderContextRows(path, file.Content, start, end, csID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"html": html})
}
