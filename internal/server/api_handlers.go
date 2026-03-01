package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
	"github.com/nikhilr/ghabricator/internal/herald"

	gh "github.com/google/go-github/v68/github"
)

// --- Task 6: Inline Comment API ---

func (s *Server) handleAPIInline(w http.ResponseWriter, r *http.Request) {
	var req APIInlineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()

	switch req.Operation {
	case "new":
		side := "LEFT"
		if req.Side == "RIGHT" || req.Side == "right" {
			side = "RIGHT"
		}

		id := nextDraftID()
		draftMu.Lock()
		draftStore[id] = &inlineDraft{
			Owner:     req.Owner,
			Repo:      req.Repo,
			Number:    req.Number,
			Path:      req.Path,
			Line:      req.Line,
			Side:      side,
			InReplyTo: req.InReplyTo,
		}
		draftMu.Unlock()

		jsonOK(w, map[string]any{
			"ok": true,
			"comment": APIInlineComment{
				ID:   id,
				Path: req.Path,
				Line: req.Line,
				Side: side,
			},
		})

	case "save":
		commentID := req.CommentID
		body := req.Body

		draftMu.Lock()
		draft, ok := draftStore[commentID]
		if ok {
			delete(draftStore, commentID)
		}
		draftMu.Unlock()

		if ok {
			var comment *ghapi.ReviewComment
			var err error
			if draft.InReplyTo > 0 {
				comment, err = ghapi.CreateReplyComment(ctx, client,
					draft.Owner, draft.Repo, draft.Number,
					body, draft.InReplyTo)
			} else {
				comment, err = ghapi.CreateReviewComment(ctx, client,
					draft.Owner, draft.Repo, draft.Number,
					body, draft.Path, draft.Line, draft.Side)
			}
			if err != nil {
				jsonError(w, err.Error(), http.StatusBadGateway)
				return
			}
			jsonOK(w, map[string]any{
				"ok": true,
				"comment": APIInlineComment{
					ID:     comment.ID,
					Author: APIUser{Login: comment.Author.Login, AvatarURL: comment.Author.AvatarURL},
					Body:   comment.Body,
					Path:   comment.Path,
					Line:   comment.Line,
					Side:   comment.Side,
				},
			})
		} else {
			comment, err := ghapi.UpdateReviewComment(ctx, client,
				req.Owner, req.Repo, commentID, body)
			if err != nil {
				jsonError(w, err.Error(), http.StatusBadGateway)
				return
			}
			jsonOK(w, map[string]any{
				"ok": true,
				"comment": APIInlineComment{
					ID:     comment.ID,
					Author: APIUser{Login: comment.Author.Login, AvatarURL: comment.Author.AvatarURL},
					Body:   comment.Body,
					Path:   comment.Path,
					Line:   comment.Line,
					Side:   comment.Side,
				},
			})
		}

	case "edit":
		comment, err := ghapi.FetchReviewComment(ctx, client, req.Owner, req.Repo, req.CommentID)
		if err != nil {
			jsonError(w, err.Error(), http.StatusBadGateway)
			return
		}
		jsonOK(w, map[string]any{
			"ok": true,
			"comment": APIInlineComment{
				ID:     comment.ID,
				Author: APIUser{Login: comment.Author.Login, AvatarURL: comment.Author.AvatarURL},
				Body:   comment.Body,
				Path:   comment.Path,
				Line:   comment.Line,
				Side:   comment.Side,
			},
		})

	case "cancel":
		draftMu.Lock()
		delete(draftStore, req.CommentID)
		draftMu.Unlock()
		jsonOK(w, map[string]bool{"ok": true})

	case "delete":
		draftMu.Lock()
		_, isDraft := draftStore[req.CommentID]
		if isDraft {
			delete(draftStore, req.CommentID)
		}
		draftMu.Unlock()

		if !isDraft {
			if err := ghapi.DeleteReviewComment(ctx, client, req.Owner, req.Repo, req.CommentID); err != nil {
				jsonError(w, err.Error(), http.StatusBadGateway)
				return
			}
		}
		jsonOK(w, map[string]bool{"ok": true})

	case "done":
		jsonOK(w, map[string]any{"ok": true, "isChecked": false})

	default:
		jsonError(w, "unknown operation: "+req.Operation, http.StatusBadRequest)
	}
}

// --- Task 7: Review/Merge/Close APIs ---

func (s *Server) handleAPIReview(w http.ResponseWriter, r *http.Request) {
	var req APIReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Owner == "" || req.Repo == "" || req.Number == 0 {
		jsonError(w, "missing owner/repo/number", http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "APPROVE", "REQUEST_CHANGES", "COMMENT":
	default:
		req.Action = "COMMENT"
	}

	client := auth.GitHubClientFromContext(r.Context())
	_, err := ghapi.SubmitReview(r.Context(), client, req.Owner, req.Repo, req.Number, req.Action, req.Body, nil)
	if err != nil {
		jsonError(w, fmt.Sprintf("submit review: %v", err), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

func (s *Server) handleAPIMerge(w http.ResponseWriter, r *http.Request) {
	var req APIMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Owner == "" || req.Repo == "" || req.Number == 0 {
		jsonError(w, "missing owner/repo/number", http.StatusBadRequest)
		return
	}

	switch req.MergeMethod {
	case "merge", "squash", "rebase":
	default:
		req.MergeMethod = "merge"
	}

	client := auth.GitHubClientFromContext(r.Context())
	opts := &gh.PullRequestOptions{MergeMethod: req.MergeMethod}
	_, _, err := client.PullRequests.Merge(r.Context(), req.Owner, req.Repo, req.Number, "", opts)
	if err != nil {
		jsonError(w, fmt.Sprintf("merge failed: %v", err), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

func (s *Server) handleAPIClose(w http.ResponseWriter, r *http.Request) {
	var req APICloseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Owner == "" || req.Repo == "" || req.Number == 0 {
		jsonError(w, "missing owner/repo/number", http.StatusBadRequest)
		return
	}

	switch req.State {
	case "closed", "open":
	default:
		jsonError(w, "invalid state", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	_, _, err := client.PullRequests.Edit(r.Context(), req.Owner, req.Repo, req.Number, &gh.PullRequest{State: &req.State})
	if err != nil {
		jsonError(w, fmt.Sprintf("update PR state: %v", err), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

// --- Reactions API ---

func (s *Server) handleAPIReaction(w http.ResponseWriter, r *http.Request) {
	var req APIReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Owner == "" || req.Repo == "" || req.CommentID == 0 || req.Content == "" {
		jsonError(w, "missing owner/repo/commentID/content", http.StatusBadRequest)
		return
	}

	switch req.Content {
	case "+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", "eyes":
	default:
		jsonError(w, "invalid reaction content", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	var err error
	if req.CommentType == "issue" {
		err = ghapi.AddIssueCommentReaction(r.Context(), client, req.Owner, req.Repo, req.CommentID, req.Content)
	} else {
		err = ghapi.AddCommentReaction(r.Context(), client, req.Owner, req.Repo, req.CommentID, req.Content)
	}
	if err != nil {
		jsonError(w, fmt.Sprintf("add reaction: %v", err), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

// --- Task 8: Repos API ---

func (s *Server) handleAPIRepos(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())

	repos, _, err := client.Repositories.List(r.Context(), "", &gh.RepositoryListOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: gh.ListOptions{PerPage: 50},
	})
	if err != nil {
		log.Printf("api repo list error: %v", err)
		jsonError(w, "failed to list repos", http.StatusBadGateway)
		return
	}

	result := make([]APIRepoSummary, 0, len(repos))
	for _, repo := range repos {
		avatarURL := ""
		if repo.Owner != nil {
			avatarURL = repo.Owner.GetAvatarURL()
		}
		result = append(result, APIRepoSummary{
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: repo.GetDescription(),
			Language:    repo.GetLanguage(),
			Stars:       repo.GetStargazersCount(),
			Forks:       repo.GetForksCount(),
			Private:     repo.GetPrivate(),
			Fork:        repo.GetFork(),
			Archived:    repo.GetArchived(),
			AvatarURL:   avatarURL,
			UpdatedAt:   timeAgo(repo.GetUpdatedAt().Time),
		})
	}
	jsonOK(w, result)
}

func repoInfoToAPI(info *ghapi.RepoInfo) APIRepoInfo {
	return APIRepoInfo{
		FullName:      info.FullName,
		Description:   info.Description,
		DefaultBranch: info.DefaultBranch,
		Private:       info.Private,
		HTMLURL:       info.HTMLURL,
		Stars:         info.Stars,
		Forks:         info.Forks,
	}
}

func (s *Server) handleAPIRepoInfo(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	client := auth.GitHubClientFromContext(r.Context())

	info, err := ghapi.FetchRepoInfo(r.Context(), client, owner, repo)
	if err != nil {
		jsonError(w, fmt.Sprintf("repo not found: %v", err), http.StatusNotFound)
		return
	}
	jsonOK(w, repoInfoToAPI(info))
}

func (s *Server) handleAPIRepoTree(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	ref := r.URL.Query().Get("ref")
	path := r.URL.Query().Get("path")
	client := auth.GitHubClientFromContext(r.Context())

	info, err := ghapi.FetchRepoInfo(r.Context(), client, owner, repo)
	if err != nil {
		jsonError(w, fmt.Sprintf("repo not found: %v", err), http.StatusNotFound)
		return
	}
	if ref == "" {
		ref = info.DefaultBranch
	}

	entries, err := ghapi.FetchRepoTree(r.Context(), client, owner, repo, ref, path)
	if err != nil {
		jsonError(w, fmt.Sprintf("tree not found: %v", err), http.StatusNotFound)
		return
	}

	apiEntries := make([]APIRepoEntry, 0, len(entries))
	for _, e := range entries {
		apiEntries = append(apiEntries, APIRepoEntry{
			Name: e.Name,
			Path: e.Path,
			Type: e.Type,
			Size: e.Size,
		})
	}

	jsonOK(w, APIRepoTreeResponse{
		Entries:  apiEntries,
		RepoInfo: repoInfoToAPI(info),
	})
}

func (s *Server) handleAPIRepoFile(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	ref := r.URL.Query().Get("ref")
	path := r.URL.Query().Get("path")
	client := auth.GitHubClientFromContext(r.Context())

	info, err := ghapi.FetchRepoInfo(r.Context(), client, owner, repo)
	if err != nil {
		jsonError(w, fmt.Sprintf("repo not found: %v", err), http.StatusNotFound)
		return
	}
	if ref == "" {
		ref = info.DefaultBranch
	}

	file, err := ghapi.FetchFileContent(r.Context(), client, owner, repo, ref, path)
	if err != nil {
		jsonError(w, fmt.Sprintf("file not found: %v", err), http.StatusNotFound)
		return
	}

	apiFile := APIRepoFile{
		Name:    file.Name,
		Path:    file.Path,
		Size:    file.Size,
		HTMLURL: file.HTMLURL,
	}

	ext := strings.ToLower(filepath.Ext(file.Name))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico", ".bmp":
		apiFile.RawURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, path)
	default:
		lines := strings.Split(file.Content, "\n")
		apiFile.Lines = diff.HighlightLines(file.Name, lines)
	}

	jsonOK(w, APIRepoFileResponse{
		File:     apiFile,
		RepoInfo: repoInfoToAPI(info),
	})
}

// --- Task 9: Paste API ---

func (s *Server) handleAPIPasteList(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())

	gists, err := ghapi.ListGists(r.Context(), client)
	if err != nil {
		log.Printf("api paste list error: %v", err)
		jsonError(w, "failed to list pastes", http.StatusBadGateway)
		return
	}

	result := make([]APIPasteSummary, 0, len(gists))
	for _, g := range gists {
		title := g.Description
		if title == "" {
			title = "(untitled)"
		}
		lang := ""
		for _, f := range g.Files {
			if f.Language != "" {
				lang = f.Language
			}
			break
		}
		result = append(result, APIPasteSummary{
			ID:        g.ID,
			Title:     title,
			Language:  lang,
			Public:    g.Public,
			CreatedAt: g.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	jsonOK(w, result)
}

func (s *Server) handleAPIPasteView(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())
	gistID := r.PathValue("id")

	gist, err := ghapi.FetchGist(r.Context(), client, gistID)
	if err != nil {
		jsonError(w, fmt.Sprintf("paste not found: %v", err), http.StatusNotFound)
		return
	}

	title := gist.Description
	if title == "" {
		title = "(untitled)"
	}

	files := make([]APIPasteFile, 0, len(gist.Files))
	for name, f := range gist.Files {
		lines := strings.Split(f.Content, "\n")
		highlighted := diff.HighlightLines(name, lines)
		files = append(files, APIPasteFile{
			Filename: name,
			Language: f.Language,
			Size:     f.Size,
			Lines:    highlighted,
		})
	}

	jsonOK(w, APIPasteDetail{
		ID:        gist.ID,
		Title:     title,
		Public:    gist.Public,
		Owner:     APIUser{Login: gist.Owner.Login, AvatarURL: gist.Owner.AvatarURL},
		HTMLURL:   gist.HTMLURL,
		Files:     files,
		CreatedAt: gist.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: gist.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (s *Server) handleAPIPasteCreate(w http.ResponseWriter, r *http.Request) {
	var req APIPasteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		jsonError(w, "content is required", http.StatusBadRequest)
		return
	}

	title := req.Title
	if title == "" {
		title = "Untitled Paste"
	}

	filename := "paste"
	if req.Language != "" {
		filename = "paste." + req.Language
	}

	client := auth.GitHubClientFromContext(r.Context())
	gist, err := ghapi.CreateGist(r.Context(), client, title, filename, req.Content, req.Public)
	if err != nil {
		jsonError(w, fmt.Sprintf("create paste failed: %v", err), http.StatusBadGateway)
		return
	}

	jsonOK(w, APIPasteCreateResponse{
		ID:  gist.ID,
		URL: "/paste/" + gist.ID,
	})
}

// --- Task 10: Herald API ---

func (s *Server) handleAPIHeraldList(w http.ResponseWriter, r *http.Request) {
	rules, err := s.herald.List()
	if err != nil {
		jsonError(w, fmt.Sprintf("load rules: %v", err), http.StatusInternalServerError)
		return
	}
	if rules == nil {
		rules = []herald.Rule{}
	}
	jsonOK(w, rules)
}

func (s *Server) handleAPIHeraldGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rule, err := s.herald.Get(id)
	if err != nil {
		jsonError(w, fmt.Sprintf("load rule: %v", err), http.StatusInternalServerError)
		return
	}
	if rule == nil {
		jsonError(w, "rule not found", http.StatusNotFound)
		return
	}
	jsonOK(w, rule)
}

func (s *Server) handleAPIHeraldSave(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())

	var rule herald.Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		jsonError(w, "bad json", http.StatusBadRequest)
		return
	}

	if rule.Name == "" {
		jsonError(w, "rule name is required", http.StatusBadRequest)
		return
	}

	rule.AuthorLogin = sess.Login

	if err := s.herald.Save(&rule); err != nil {
		jsonError(w, fmt.Sprintf("save rule: %v", err), http.StatusInternalServerError)
		return
	}
	jsonOK(w, rule)
}

func (s *Server) handleAPIHeraldDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.herald.Delete(id); err != nil {
		jsonError(w, fmt.Sprintf("delete rule: %v", err), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

// --- Task 11: Search API ---

func (s *Server) handleAPISearch(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")

	if query == "" {
		jsonOK(w, APISearchResponse{})
		return
	}

	if searchType == "" {
		searchType = "pr"
	}

	var resp APISearchResponse

	switch searchType {
	case "pr":
		prs := s.searchPRs(r, client, "is:pr "+query)
		apiPRs := make([]APISearchPR, 0, len(prs))
		for _, pr := range prs {
			apiPRs = append(apiPRs, APISearchPR{
				Number:    pr.Number,
				Title:     pr.Title,
				Repo:      pr.Repo,
				Author:    pr.Author,
				AvatarURL: pr.AvatarURL,
				UpdatedAt: pr.UpdatedAt,
				Draft:     pr.Draft,
				URL:       pr.URL,
			})
		}
		resp.PRs = apiPRs

	case "code":
		result, _, err := client.Search.Code(r.Context(), query, &gh.SearchOptions{
			ListOptions: gh.ListOptions{PerPage: 25},
		})
		if err != nil {
			log.Printf("api search code error: %v", err)
			jsonError(w, "code search failed", http.StatusBadGateway)
			return
		}
		codes := make([]APISearchCodeResult, 0, len(result.CodeResults))
		for _, cr := range result.CodeResults {
			fragment := ""
			if len(cr.TextMatches) > 0 {
				fragment = cr.TextMatches[0].GetFragment()
			}
			codes = append(codes, APISearchCodeResult{
				Repo:     cr.GetRepository().GetFullName(),
				Path:     cr.GetPath(),
				Fragment: fragment,
			})
		}
		resp.Code = codes

	case "repo":
		result, _, err := client.Search.Repositories(r.Context(), query, &gh.SearchOptions{
			Sort:        "stars",
			Order:       "desc",
			ListOptions: gh.ListOptions{PerPage: 25},
		})
		if err != nil {
			log.Printf("api search repos error: %v", err)
			jsonError(w, "repo search failed", http.StatusBadGateway)
			return
		}
		repos := make([]APISearchRepoResult, 0, len(result.Repositories))
		for _, repo := range result.Repositories {
			avatarURL := ""
			if repo.Owner != nil {
				avatarURL = repo.Owner.GetAvatarURL()
			}
			repos = append(repos, APISearchRepoResult{
				FullName:    repo.GetFullName(),
				Description: repo.GetDescription(),
				Stars:       repo.GetStargazersCount(),
				Language:    repo.GetLanguage(),
				AvatarURL:   avatarURL,
			})
		}
		resp.Repos = repos
	}

	jsonOK(w, resp)
}
