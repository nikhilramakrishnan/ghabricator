package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nikhilr/ghabricator/internal/auth"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
)

func (s *Server) handleReview(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()
	isJSON := strings.Contains(r.Header.Get("Content-Type"), "application/json")

	var owner, repo, event, body string
	var number int
	var inlineComments []ghapi.InlineCommentRequest

	if isJSON {
		var req struct {
			Owner    string                      `json:"owner"`
			Repo     string                      `json:"repo"`
			Number   int                         `json:"number"`
			Action   string                      `json:"action"`
			Body     string                      `json:"body"`
			Comments []ghapi.InlineCommentRequest `json:"comments,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad json"})
			return
		}
		owner, repo, number = req.Owner, req.Repo, req.Number
		event, body = req.Action, req.Body
		inlineComments = req.Comments
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		}
		owner = r.FormValue("owner")
		repo = r.FormValue("repo")
		number, _ = strconv.Atoi(r.FormValue("number"))
		event = r.FormValue("action")
		body = r.FormValue("body")
	}

	if owner == "" || repo == "" || number == 0 {
		if isJSON {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing owner/repo/number"})
		} else {
			http.Error(w, "missing owner/repo/number", http.StatusBadRequest)
		}
		return
	}

	switch event {
	case "APPROVE", "REQUEST_CHANGES", "COMMENT":
	default:
		event = "COMMENT"
	}

	_, err := ghapi.SubmitReview(ctx, client, owner, repo, number, event, body, inlineComments)
	if err != nil {
		errMsg := fmt.Sprintf("submit review: %v", err)
		if isJSON {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": errMsg})
		} else if r.FormValue("__wflow__") == "true" {
			writeJavelinPayload(w, http.StatusOK, map[string]any{"error": errMsg})
		} else {
			http.Error(w, errMsg, http.StatusBadGateway)
		}
		return
	}

	redirect := fmt.Sprintf("/pr/%s/%s/%d", owner, repo, number)
	if isJSON {
		writeJSON(w, http.StatusOK, map[string]string{"redirect": redirect})
	} else if r.FormValue("__wflow__") == "true" {
		// Javelin Workflow expects redirect in payload.
		writeJavelinPayload(w, http.StatusOK, map[string]any{"redirect": redirect})
	} else {
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}
