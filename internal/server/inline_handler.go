package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
)

// Draft inline comment stored in memory between op=new and op=save.
type inlineDraft struct {
	Owner  string
	Repo   string
	Number int
	Path   string
	Line   int
	Side   string // LEFT or RIGHT
}

var (
	draftMu    sync.Mutex
	draftSeq   atomic.Int64
	draftStore = make(map[int64]*inlineDraft)
)

func nextDraftID() int64 {
	return draftSeq.Add(1)
}

// writeJavelinPayload writes a Javelin-compatible AJAX response:
//
//	for (;;);{"payload":{...}}
func writeJavelinPayload(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]any{"payload": payload}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(w, "for (;;);%s", data)
}

// inlineContentState mirrors DiffInlineContentState wire format.
func emptyContentState() map[string]any {
	return map[string]any{
		"text":           "",
		"suggestionText": nil,
		"hasSuggestion":  false,
	}
}

func (s *Server) handleInline(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJavelinPayload(w, http.StatusBadRequest, map[string]any{"error": "bad form data"})
		return
	}

	op := r.FormValue("op")
	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()

	// owner/repo/number come from query params (set in inlineURI).
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	number, _ := strconv.Atoi(r.URL.Query().Get("number"))

	switch op {
	case "new":
		// changesetID is the file path (set in renderer data-meta left/right).
		path := r.FormValue("changesetID")
		lineNum, _ := strconv.Atoi(r.FormValue("number"))
		onRight := r.FormValue("on_right") == "1"

		side := "LEFT"
		if onRight {
			side = "RIGHT"
		}

		// Store draft for later save.
		id := nextDraftID()
		draftMu.Lock()
		draftStore[id] = &inlineDraft{
			Owner:  owner,
			Repo:   repo,
			Number: number,
			Path:   path,
			Line:   lineNum,
			Side:   side,
		}
		draftMu.Unlock()

		view := renderInlineEditor(id, lineNum, onRight, "")
		writeJavelinPayload(w, http.StatusOK, map[string]any{
			"view": view,
			"inline": map[string]any{
				"id": id,
				"state": map[string]any{
					"initial":   nil,
					"committed": nil,
					"active":    emptyContentState(),
				},
				"canSuggestEdit": false,
			},
		})

	case "save":
		commentID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		body := r.FormValue("text")

		draftMu.Lock()
		draft, ok := draftStore[commentID]
		if ok {
			delete(draftStore, commentID)
		}
		draftMu.Unlock()

		if ok {
			// New comment — create on GitHub.
			comment, err := ghapi.CreateReviewComment(ctx, client,
				draft.Owner, draft.Repo, draft.Number,
				body, draft.Path, draft.Line, draft.Side)
			if err != nil {
				writeJavelinPayload(w, http.StatusOK, map[string]any{
					"error": err.Error(),
				})
				return
			}
			view := renderInlineCommentHTML(comment)
			writeJavelinPayload(w, http.StatusOK, map[string]any{
				"view": view,
			})
		} else {
			// Existing comment — update on GitHub.
			comment, err := ghapi.UpdateReviewComment(ctx, client,
				owner, repo, commentID, body)
			if err != nil {
				writeJavelinPayload(w, http.StatusOK, map[string]any{
					"error": err.Error(),
				})
				return
			}
			view := renderInlineCommentHTML(comment)
			writeJavelinPayload(w, http.StatusOK, map[string]any{
				"view": view,
			})
		}

	case "edit":
		commentID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		// Fetch the existing comment to pre-fill the editor.
		comment, err := ghapi.FetchReviewComment(ctx, client, owner, repo, commentID)
		if err != nil {
			writeJavelinPayload(w, http.StatusOK, map[string]any{
				"error": err.Error(),
			})
			return
		}
		view := renderInlineEditor(commentID, comment.Line, comment.Side == "RIGHT", comment.Body)
		state := emptyContentState()
		state["text"] = comment.Body
		writeJavelinPayload(w, http.StatusOK, map[string]any{
			"view": view,
			"inline": map[string]any{
				"id": commentID,
				"state": map[string]any{
					"initial":   state,
					"committed": state,
					"active":    state,
				},
				"canSuggestEdit": false,
			},
		})

	case "cancel":
		commentID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		// Clean up draft if it exists.
		draftMu.Lock()
		delete(draftStore, commentID)
		draftMu.Unlock()
		writeJavelinPayload(w, http.StatusOK, map[string]any{})

	case "delete":
		commentID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

		// If it's a draft, just remove it.
		draftMu.Lock()
		_, isDraft := draftStore[commentID]
		if isDraft {
			delete(draftStore, commentID)
		}
		draftMu.Unlock()

		if !isDraft {
			err := ghapi.DeleteReviewComment(ctx, client, owner, repo, commentID)
			if err != nil {
				writeJavelinPayload(w, http.StatusOK, map[string]any{
					"error": err.Error(),
				})
				return
			}
		}
		writeJavelinPayload(w, http.StatusOK, map[string]any{})

	case "draft":
		// Auto-save drafts — no-op for us since we don't persist drafts.
		writeJavelinPayload(w, http.StatusOK, map[string]any{})

	case "done":
		// Toggle "done" state — not implemented, just ack.
		writeJavelinPayload(w, http.StatusOK, map[string]any{
			"isChecked":  false,
			"draftState": false,
		})

	default:
		writeJavelinPayload(w, http.StatusBadRequest, map[string]any{
			"error": "unknown op: " + op,
		})
	}
}

func renderInlineEditor(id int64, lineNum int, isNew bool, existingText string) string {
	var b strings.Builder
	side := "new"
	if !isNew {
		side = "old"
	}
	fmt.Fprintf(&b, `<tr class="inline inline-editor" data-sigil="inline-row" data-line="%d" data-side="%s">`, lineNum, side)
	b.WriteString(`<td colspan="6">`)
	b.WriteString(`<div class="differential-inline-comment-edit" data-sigil="differential-inline-comment">`)
	// Entire editor wrapped in form so Save button triggers submit.
	b.WriteString(`<form data-sigil="inline-edit-form">`)
	b.WriteString(`<div class="differential-inline-comment-edit-title">New Inline Comment</div>`)
	b.WriteString(`<div class="differential-inline-comment-edit-body">`)
	fmt.Fprintf(&b, `<div class="aphront-form-input"><textarea data-sigil="inline-content-text" class="remarkup-assist-textarea PhabricatorMonospaced" rows="4">%s</textarea></div>`,
		template.HTMLEscapeString(existingText))
	b.WriteString(`</div>`)
	b.WriteString(`<div class="differential-inline-comment-edit-buttons" data-sigil="inline-edit-buttons">`)
	b.WriteString(`<button type="submit" class="phui-button-view button-green">Save</button>`)
	b.WriteString(` <button type="button" class="phui-button-view button-grey" data-sigil="inline-edit-cancel">Cancel</button>`)
	b.WriteString(`</div>`)
	b.WriteString(`</form>`)
	b.WriteString(`</div>`)
	b.WriteString(`</td></tr>`)
	return b.String()
}

func renderInlineCommentHTML(c *ghapi.ReviewComment) string {
	var b strings.Builder

	ic := diff.InlineComment{
		ID:        c.ID,
		Author:    c.Author.Login,
		AvatarURL: c.Author.AvatarURL,
		Body:      c.Body,
		Path:      c.Path,
		Line:      c.Line,
		Side:      c.Side,
	}
	metaJSON, _ := json.Marshal(diff.InlineCommentMeta(ic))

	b.WriteString(`<tr class="inline" data-sigil="inline-row">`)
	b.WriteString(`<td colspan="6">`)
	fmt.Fprintf(&b, `<div class="differential-inline-comment inline-comment-element mood-inline-comment" data-sigil="differential-inline-comment" data-meta='%s'>`,
		template.HTMLEscapeString(string(metaJSON)))

	b.WriteString(`<div class="mood-inline-header">`)
	if c.Author.AvatarURL != "" {
		fmt.Fprintf(&b, `<img src="%s" style="width:20px;height:20px;border-radius:3px">`, template.HTMLEscapeString(c.Author.AvatarURL))
	}
	fmt.Fprintf(&b, `<strong>%s</strong>`, template.HTMLEscapeString(c.Author.Login))
	b.WriteString(`</div>`)

	b.WriteString(`<div class="mood-inline-body">`)
	b.WriteString(`<div class="phabricator-remarkup">`)
	b.WriteString(template.HTMLEscapeString(c.Body))
	b.WriteString(`</div></div>`)

	b.WriteString(`<div class="mood-inline-actions">`)
	b.WriteString(`<a class="inline-action-reply"><span class="phui-icon-view phui-font-fa fa-reply"></span> Reply</a>`)
	b.WriteString(`<a class="inline-action-done"><span class="phui-icon-view phui-font-fa fa-check"></span> Done</a>`)
	b.WriteString(`</div>`)

	b.WriteString(`</div>`)
	b.WriteString(`</td></tr>`)
	return b.String()
}

// writeJSON writes a plain JSON response (for non-Javelin callers).
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
