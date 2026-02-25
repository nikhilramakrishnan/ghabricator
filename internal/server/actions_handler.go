package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/nikhilr/ghabricator/internal/auth"
	gh "github.com/google/go-github/v68/github"
)

func (s *Server) handleMerge(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	owner := r.FormValue("owner")
	repo := r.FormValue("repo")
	number, _ := strconv.Atoi(r.FormValue("number"))
	method := r.FormValue("merge_method")

	if owner == "" || repo == "" || number == 0 {
		http.Error(w, "missing owner/repo/number", http.StatusBadRequest)
		return
	}

	switch method {
	case "merge", "squash", "rebase":
	default:
		method = "merge"
	}

	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()

	opts := &gh.PullRequestOptions{MergeMethod: method}
	_, _, err := client.PullRequests.Merge(ctx, owner, repo, number, "", opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("merge failed: %v", err), http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/pr/%s/%s/%d",
		template.URLQueryEscaper(owner),
		template.URLQueryEscaper(repo),
		number), http.StatusSeeOther)
}

func (s *Server) handleClose(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	owner := r.FormValue("owner")
	repo := r.FormValue("repo")
	number, _ := strconv.Atoi(r.FormValue("number"))
	state := r.FormValue("state")

	if owner == "" || repo == "" || number == 0 {
		http.Error(w, "missing owner/repo/number", http.StatusBadRequest)
		return
	}

	switch state {
	case "closed", "open":
	default:
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()

	_, _, err := client.PullRequests.Edit(ctx, owner, repo, number, &gh.PullRequest{State: &state})
	if err != nil {
		http.Error(w, fmt.Sprintf("update PR state: %v", err), http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/pr/%s/%s/%d",
		template.URLQueryEscaper(owner),
		template.URLQueryEscaper(repo),
		number), http.StatusSeeOther)
}
