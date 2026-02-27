package server

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
	"github.com/nikhilr/ghabricator/internal/herald"
	"github.com/nikhilr/ghabricator/internal/remarkup"
)

// timelineEvent is a unified event for the PR activity timeline.
type timelineEvent struct {
	Author    ghapi.User
	Action    string // "approved this revision", "requested changes to this revision", "added a comment", etc.
	Body      string
	CreatedAt time.Time
	IconClass string // FA icon class
	IconColor string // CSS color class
}

func buildTimeline(pr *ghapi.PullRequest, reviews []ghapi.Review, issueComments []ghapi.IssueComment) []timelineEvent {
	var events []timelineEvent

	// PR creation event (body shown in curtain Summary, not here).
	events = append(events, timelineEvent{
		Author:    pr.Author,
		Action:    "created this revision",
		CreatedAt: pr.CreatedAt,
		IconClass: "fa-plus",
		IconColor: "blue",
	})

	// Reviews.
	for _, r := range reviews {
		var action, icon, color string
		switch r.State {
		case "APPROVED":
			action = "accepted this revision"
			icon = "fa-check-circle"
			color = "green"
		case "CHANGES_REQUESTED":
			action = "requested changes to this revision"
			icon = "fa-times-circle"
			color = "red"
		case "COMMENTED":
			action = "added a comment"
			icon = "fa-comment"
			color = "blue"
		case "DISMISSED":
			action = "had their review dismissed"
			icon = "fa-ban"
			color = "grey"
		default:
			action = "reviewed"
			icon = "fa-comment"
			color = "blue"
		}
		events = append(events, timelineEvent{
			Author:    r.Author,
			Action:    action,
			Body:      r.Body,
			CreatedAt: r.CreatedAt,
			IconClass: icon,
			IconColor: color,
		})
	}

	// Issue comments.
	for _, c := range issueComments {
		events = append(events, timelineEvent{
			Author:    c.Author,
			Action:    "added a comment",
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			IconClass: "fa-comment",
			IconColor: "blue",
		})
	}

	// Merged event.
	if pr.Merged {
		events = append(events, timelineEvent{
			Author:    pr.Author,
			Action:    "closed this revision",
			CreatedAt: pr.UpdatedAt,
			IconClass: "fa-check",
			IconColor: "violet",
		})
	} else if pr.State == "closed" {
		events = append(events, timelineEvent{
			Author:    pr.Author,
			Action:    "abandoned this revision",
			CreatedAt: pr.UpdatedAt,
			IconClass: "fa-ban",
			IconColor: "red",
		})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	return events
}

// --- JSON API handler ---

func (s *Server) handleAPIPR(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	numberStr := r.PathValue("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		jsonError(w, "invalid PR number", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	ctx := r.Context()

	// Parallel fetch — same as HTML handler.
	var (
		pr            *ghapi.PullRequest
		rawDiff       string
		comments      []ghapi.ReviewComment
		reviews       []ghapi.Review
		issueComments []ghapi.IssueComment
		commits       []ghapi.PRCommit
		prErr, diffErr, commentsErr, reviewsErr, issueCommentsErr, commitsErr error
	)

	var wg sync.WaitGroup
	wg.Add(6)
	go func() { defer wg.Done(); pr, prErr = ghapi.FetchPR(ctx, client, owner, repo, number) }()
	go func() { defer wg.Done(); rawDiff, diffErr = ghapi.FetchDiff(ctx, client, owner, repo, number) }()
	go func() {
		defer wg.Done()
		comments, commentsErr = ghapi.FetchReviewComments(ctx, client, owner, repo, number)
	}()
	go func() { defer wg.Done(); reviews, reviewsErr = ghapi.FetchReviews(ctx, client, owner, repo, number) }()
	go func() {
		defer wg.Done()
		issueComments, issueCommentsErr = ghapi.FetchIssueComments(ctx, client, owner, repo, number)
	}()
	go func() {
		defer wg.Done()
		commits, commitsErr = ghapi.FetchPRCommits(ctx, client, owner, repo, number)
	}()
	wg.Wait()

	if prErr != nil {
		jsonError(w, fmt.Sprintf("could not load PR: %v", prErr), http.StatusBadGateway)
		return
	}
	if diffErr != nil {
		jsonError(w, fmt.Sprintf("could not load diff: %v", diffErr), http.StatusBadGateway)
		return
	}

	// Fetch check runs (need head SHA).
	checkRuns, checksErr := ghapi.FetchCheckRuns(ctx, client, owner, repo, pr.Head.SHA)
	if checksErr != nil {
		checkRuns = nil
	}

	// Non-fatal errors.
	if commentsErr != nil {
		comments = nil
	}
	if reviewsErr != nil {
		reviews = nil
	}
	if issueCommentsErr != nil {
		issueComments = nil
	}
	if commitsErr != nil {
		commits = nil
	}

	// Parse diff.
	changesets, err := diff.ParseDiff(rawDiff)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not parse diff: %v", err), http.StatusInternalServerError)
		return
	}

	// Index comments by path.
	commentsByPath := make(map[string][]ghapi.ReviewComment)
	for _, c := range comments {
		commentsByPath[c.Path] = append(commentsByPath[c.Path], c)
	}

	// Build changesets with diff rows.
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

	// Build API comments by path.
	apiCommentsByPath := make(map[string][]APIReviewComment, len(commentsByPath))
	for path, cmts := range commentsByPath {
		apiCmts := make([]APIReviewComment, 0, len(cmts))
		for _, c := range cmts {
			apiCmts = append(apiCmts, APIReviewComment{
				ID:        c.ID,
				Author:    APIUser{Login: c.Author.Login, AvatarURL: c.Author.AvatarURL},
				Body:      c.Body,
				Path:      c.Path,
				Line:      c.Line,
				Side:      c.Side,
				CreatedAt: c.CreatedAt,
				InReplyTo: c.InReplyTo,
			})
		}
		apiCommentsByPath[path] = apiCmts
	}

	// Build reviews.
	apiReviews := make([]APIReview, 0, len(reviews))
	for _, rv := range reviews {
		apiReviews = append(apiReviews, APIReview{
			ID:        rv.ID,
			Author:    APIUser{Login: rv.Author.Login, AvatarURL: rv.Author.AvatarURL},
			State:     rv.State,
			Body:      rv.Body,
			CreatedAt: rv.CreatedAt,
		})
	}

	// Build issue comments.
	apiIssueComments := make([]APIIssueComment, 0, len(issueComments))
	for _, ic := range issueComments {
		apiIssueComments = append(apiIssueComments, APIIssueComment{
			ID:        ic.ID,
			Author:    APIUser{Login: ic.Author.Login, AvatarURL: ic.Author.AvatarURL},
			Body:      ic.Body,
			CreatedAt: ic.CreatedAt,
		})
	}

	// Build check runs.
	apiCheckRuns := make([]APICheckRun, 0, len(checkRuns))
	for _, cr := range checkRuns {
		apiCheckRuns = append(apiCheckRuns, APICheckRun{
			Name:        cr.Name,
			Status:      cr.Status,
			Conclusion:  cr.Conclusion,
			DetailsURL:  cr.DetailsURL,
			AppName:     cr.AppName,
			StartedAt:   cr.StartedAt,
			CompletedAt: cr.CompletedAt,
		})
	}

	// Build timeline.
	events := buildTimeline(pr, reviews, issueComments)
	apiTimeline := make([]APITimelineEvent, 0, len(events))
	for _, ev := range events {
		apiTimeline = append(apiTimeline, APITimelineEvent{
			Author:    APIUser{Login: ev.Author.Login, AvatarURL: ev.Author.AvatarURL},
			Action:    ev.Action,
			Body:      ev.Body,
			CreatedAt: ev.CreatedAt,
			IconClass: ev.IconClass,
			IconColor: ev.IconColor,
		})
	}

	// Herald evaluation.
	var apiHeraldMatches []APIHeraldMatch
	if rules, heraldErr := s.herald.List(); heraldErr == nil && len(rules) > 0 {
		var changedFiles []string
		for _, cs := range changesets {
			changedFiles = append(changedFiles, cs.DisplayPath())
		}
		var labels []string
		for _, l := range pr.Labels {
			labels = append(labels, l.Name)
		}
		prCtx := &herald.PRContext{
			Author:       pr.Author.Login,
			Title:        pr.Title,
			Labels:       labels,
			BaseBranch:   pr.Base.Ref,
			ChangedFiles: changedFiles,
		}
		matches := herald.Evaluate(rules, prCtx)
		for _, m := range matches {
			am := APIHeraldMatch{
				RuleID:   m.Rule.ID,
				RuleName: m.Rule.Name,
			}
			for _, a := range m.Actions {
				am.Actions = append(am.Actions, APIHeraldAction{
					Type:  string(a.Type),
					Value: a.Value,
				})
			}
			apiHeraldMatches = append(apiHeraldMatches, am)
		}
	}

	// Build commits.
	apiCommits := make([]APICommit, 0, len(commits))
	for _, c := range commits {
		apiCommits = append(apiCommits, APICommit{
			SHA:     c.SHA,
			Message: c.Message,
			Author:  APIUser{Login: c.Author.Login, AvatarURL: c.Author.AvatarURL},
			Date:    c.Date.Format("2006-01-02T15:04:05Z"),
		})
	}

	// Build PR detail.
	apiLabels := make([]APILabel, 0, len(pr.Labels))
	for _, l := range pr.Labels {
		apiLabels = append(apiLabels, APILabel{Name: l.Name, Color: l.Color})
	}
	apiReviewers := make([]APIUser, 0, len(pr.Reviewers))
	for _, u := range pr.Reviewers {
		apiReviewers = append(apiReviewers, APIUser{Login: u.Login, AvatarURL: u.AvatarURL})
	}

	resp := APIPRDetailResponse{
		PR: APIPRDetail{
			Number:       pr.Number,
			Title:        pr.Title,
			Body:         remarkup.Render(pr.Body),
			State:        pr.State,
			Draft:        pr.Draft,
			Merged:       pr.Merged,
			Author:       APIUser{Login: pr.Author.Login, AvatarURL: pr.Author.AvatarURL},
			CreatedAt:    pr.CreatedAt,
			UpdatedAt:    pr.UpdatedAt,
			Labels:       apiLabels,
			Reviewers:    apiReviewers,
			Head:         APIRef{Ref: pr.Head.Ref, SHA: pr.Head.SHA, Repo: pr.Head.Repo},
			Base:         APIRef{Ref: pr.Base.Ref, SHA: pr.Base.SHA, Repo: pr.Base.Repo},
			Additions:    pr.Additions,
			Deletions:    pr.Deletions,
			ChangedFiles: pr.ChangedFiles,
		},
		Changesets:     apiChangesets,
		CommentsByPath: apiCommentsByPath,
		Reviews:        apiReviews,
		IssueComments:  apiIssueComments,
		CheckRuns:      apiCheckRuns,
		Timeline:       apiTimeline,
		HeraldMatches:  apiHeraldMatches,
		Commits:        apiCommits,
	}

	jsonOK(w, resp)
}
