package server

import (
	"fmt"
	"net/http"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
)

func (s *Server) handleAPICompare(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	base := r.URL.Query().Get("base")
	head := r.URL.Query().Get("head")

	if base == "" || head == "" {
		jsonError(w, "base and head query params required", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()

	rawDiff, err := ghapi.FetchCompare(ctx, client, owner, repo, base, head)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not fetch compare diff: %v", err), http.StatusBadGateway)
		return
	}

	changesets, err := diff.ParseDiff(rawDiff)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not parse diff: %v", err), http.StatusInternalServerError)
		return
	}

	apiChangesets := make([]APIChangeset, 0, len(changesets))
	for _, cs := range changesets {
		rows := diff.BuildDiffRows(cs)
		apiRows := make([]APIDiffRow, 0, len(rows))
		for _, row := range rows {
			apiRows = append(apiRows, APIDiffRow{
				OldNum:     row.OldNum,
				NewNum:     row.NewNum,
				OldClass:   row.OldClass,
				NewClass:   row.NewClass,
				OldContent: string(row.OldContent),
				NewContent: string(row.NewContent),
				IsContext:  row.IsContext,
			})
		}
		apiChangesets = append(apiChangesets, APIChangeset{
			ID:           cs.ID,
			OldName:      cs.OldName,
			NewName:      cs.NewName,
			DisplayPath:  cs.DisplayPath(),
			LinesAdded:   cs.LinesAdded,
			LinesRemoved: cs.LinesRemoved,
			IsNew:        cs.IsNew,
			IsDeleted:    cs.IsDeleted,
			IsRenamed:    cs.IsRenamed,
			IsBinary:     cs.IsBinary,
			Rows:         apiRows,
		})
	}

	jsonOK(w, map[string]any{
		"changesets": apiChangesets,
	})
}
