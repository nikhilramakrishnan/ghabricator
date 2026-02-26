package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
	"github.com/nikhilr/ghabricator/internal/herald"
	"github.com/nikhilr/ghabricator/internal/remarkup"
	"github.com/nikhilr/ghabricator/internal/templates"
)

func (s *Server) handlePR(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	numberStr := r.PathValue("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		s.renderError(w, r, "Invalid Request", "The PR number provided is not valid.", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())
	sess := auth.SessionFromContext(r.Context())
	ctx := r.Context()

	// Fetch all data in parallel.
	var (
		pr            *ghapi.PullRequest
		rawDiff       string
		comments      []ghapi.ReviewComment
		reviews       []ghapi.Review
		issueComments []ghapi.IssueComment
		prErr, diffErr, commentsErr, reviewsErr, issueCommentsErr error
	)

	var wg sync.WaitGroup
	wg.Add(5)
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
	wg.Wait()

	if prErr != nil {
		s.renderError(w, r, "Fetch Error", fmt.Sprintf("Could not load pull request: %v", prErr), http.StatusBadGateway)
		return
	}
	if diffErr != nil {
		s.renderError(w, r, "Fetch Error", fmt.Sprintf("Could not load diff: %v", diffErr), http.StatusBadGateway)
		return
	}
	if commentsErr != nil {
		comments = nil // non-fatal
	}
	if reviewsErr != nil {
		reviews = nil // non-fatal
	}
	if issueCommentsErr != nil {
		issueComments = nil // non-fatal
	}

	// Parse diff.
	changesets, err := diff.ParseDiff(rawDiff)
	if err != nil {
		s.renderError(w, r, "Parse Error", fmt.Sprintf("Could not parse diff: %v", err), http.StatusInternalServerError)
		return
	}

	// Index comments by file path.
	commentsByPath := make(map[string][]ghapi.ReviewComment)
	for _, c := range comments {
		commentsByPath[c.Path] = append(commentsByPath[c.Path], c)
	}

	// Build file tree sidebar.
	fileTree := diff.RenderFileTree(changesets)

	// Build content: summary on top (full-width, Phabricator-style), then changesets.
	var content strings.Builder

	// Summary (PR description) — full-width above the diff, like Phabricator.
	if strings.TrimSpace(pr.Body) != "" {
		content.WriteString(`<div class="phui-box phui-box-border phui-object-box mlt mlr">`)
		content.WriteString(`<div class="phui-header-shell"><div class="phui-header-view"><h1 class="phui-header-header">`)
		content.WriteString(`<span class="phui-header-icon phui-icon-view phui-font-fa fa-file-text-o"></span>Summary</h1></div></div>`)
		content.WriteString(`<div style="padding:10px 12px"><div class="phabricator-remarkup" style="font-size:13px;line-height:1.5;overflow-wrap:break-word;word-break:break-word;">`)
		content.WriteString(remarkup.Render(pr.Body))
		content.WriteString(`</div></div></div>`)
	}

	// Also build Javelin metadata block for changeset data-meta pointers.

	changesetViewIDs := make([]string, 0, len(changesets))
	metaBlock := make(map[string]any) // Javelin metadata block 0
	metaIdx := 0

	for _, cs := range changesets {
		viewID := fmt.Sprintf("diff-C%d", cs.ID)
		changesetViewIDs = append(changesetViewIDs, viewID)

		// Store changeset meta in the Javelin data block.
		metaKey := fmt.Sprintf("%d", metaIdx)
		metaRef := fmt.Sprintf("0_%d", metaIdx)
		metaBlock[metaKey] = diff.BuildChangesetMeta(cs)
		metaIdx++

		// Convert GitHub review comments to diff inline comments.
		path := cs.DisplayPath()
		var inlineComments []diff.InlineComment
		if fileComments, ok := commentsByPath[path]; ok {
			for _, c := range fileComments {
				inlineComments = append(inlineComments, diff.InlineComment{
					ID:        c.ID,
					Author:    c.Author.Login,
					AvatarURL: c.Author.AvatarURL,
					Body:      c.Body,
					Path:      path,
					Line:      c.Line,
					Side:      c.Side,
				})
			}
		}

		content.WriteString(string(diff.RenderChangeset(cs, metaRef, inlineComments)))
	}

	// Render timeline between diff and review form.
	renderTimeline(&content, pr, reviews, issueComments)

	// Render review form.
	renderReviewForm(&content, owner, repo, number, sess)

	// Build curtain (sidebar) — moodboard order: Reviewers, Labels, Herald, Properties, Actions.
	var curtain strings.Builder
	renderCurtainReviewersLabels(&curtain, pr, reviews)

	// Herald (between Labels and Properties, matching moodboard)
	if rules, err := s.herald.List(); err == nil && len(rules) > 0 {
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
		renderHeraldCurtainPanel(&curtain, matches)
	} else {
		renderHeraldCurtainPanel(&curtain, nil)
	}

	renderCurtainPropertiesActions(&curtain, pr, owner, repo, number)

	// Status badge.
	statusText, statusColor := prStatusBadge(pr)
	headerTitle := template.HTML(fmt.Sprintf(
		`<span class="phui-tag-view phui-tag-shade-%s mrs">`+
			`<span class="phui-tag-core">%s</span></span> `+
			`<span class="phui-header-subheader" style="color:rgba(55,55,55,.6);font-weight:normal">D%d</span> %s`,
		statusColor, statusText, number, template.HTMLEscapeString(pr.Title)))

	// Javelin metadata block — loaded before behaviors via kind="merge".
	mergeJSON, _ := json.Marshal(map[string]any{
		"block": 0,
		"data":  metaBlock,
	})

	// Javelin init data for differential behavior.
	inlineURI := fmt.Sprintf("/api/inline?owner=%s&repo=%s&number=%d", owner, repo, number)
	behaviorData := map[string]any{
		"changesetViewIDs": changesetViewIDs,
		"inlineURI":        inlineURI,
		"inlineListURI":    inlineURI,
		"isStandalone":     false,
		"pht":              map[string]string{},
	}
	behaviorsJSON, _ := json.Marshal(map[string]any{
		"differential-populate": []any{behaviorData},
	})

	theme := templates.ThemeFromRequest(r)
	extraCSS := []string{"/res/pkg/differential.pkg.css"}
	if theme == "dark" {
		extraCSS = []string{"/res/pkg/dark/differential.pkg.css"}
	}

	// Context expansion script — wires up "show more lines" clicks.
	contextScript := template.JS(fmt.Sprintf(`(function(){
  var owner=%q, repo=%q, ref=%q;
  document.addEventListener('click', function(e) {
    var td = e.target.closest('[data-action="context-expand"]');
    if (!td) return;
    e.stopPropagation();
    var path = td.getAttribute('data-path');
    var start = td.getAttribute('data-start');
    var end = td.getAttribute('data-end');
    var cs = td.getAttribute('data-cs');
    td.textContent = 'Loading\u2026';
    td.style.cursor = 'wait';
    fetch('/api/context?owner='+encodeURIComponent(owner)+'&repo='+encodeURIComponent(repo)+
      '&ref='+encodeURIComponent(ref)+'&path='+encodeURIComponent(path)+
      '&start='+start+'&end='+end+'&cs='+cs)
    .then(function(r){return r.json()})
    .then(function(data){
      var tr = td.closest('tr');
      if (tr && data.html) {
        tr.insertAdjacentHTML('afterend', data.html);
        tr.remove();
      }
    })
    .catch(function(){td.textContent='Failed to load context';td.style.cursor='pointer'});
  }, true);
})()`, owner, repo, pr.Head.Ref))

	templates.RenderPage(w, templates.PageData{
		Title:         fmt.Sprintf("D%d: %s", number, pr.Title),
		Theme:         theme,
		HeaderTitle:   headerTitle,
		HeaderIcon:    "fa-cog",
		BodyClass:     "phui-workboard-color",
		Content:       template.HTML(content.String()),
		Curtain:       template.HTML(curtain.String()),
		FileTree:      fileTree,
		ExtraCSS:      extraCSS,
		ExtraJS:       []string{"/res/pkg/differential.pkg.js"},
		InlineScript:  contextScript,
		NavActive:     "revisions",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: owner + "/" + repo, Href: "/dashboard"},
			{Name: fmt.Sprintf("#%d", number)},
		},
		JavelinData: []templates.JavelinInit{
			{Kind: "merge", Data: template.HTML(mergeJSON)},
			{Kind: "behaviors", Data: template.HTML(behaviorsJSON)},
		},
	})
}

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

func renderTimeline(b *strings.Builder, pr *ghapi.PullRequest, reviews []ghapi.Review, issueComments []ghapi.IssueComment) {
	events := buildTimeline(pr, reviews, issueComments)
	if len(events) == 0 {
		return
	}

	b.WriteString(`<div class="phui-box phui-box-border phui-object-box" style="padding:10px 12px; margin:12px 0 0 0;">`)
	b.WriteString(`<div class="phui-timeline-view">`)
	for i, ev := range events {
		renderTimelineEvent(b, ev, i == len(events)-1)
	}
	b.WriteString(`</div>`)
	b.WriteString(`</div>`)
}

func renderTimelineEvent(b *strings.Builder, ev timelineEvent, isLast bool) {
	isMajor := ev.Body != ""
	esc := template.HTMLEscapeString

	// Moodboard-style: simple flex row with colored circle icon
	borderStyle := ` border-bottom:1px solid #e3e4e8;`
	if isLast {
		borderStyle = ""
	}
	fmt.Fprintf(b, `<div class="phui-timeline-event-view" style="display:flex; gap:10px; padding:10px 0;%s">`, borderStyle)

	// Icon circle or avatar
	if isMajor && ev.Author.AvatarURL != "" {
		fmt.Fprintf(b, `<img src="%s" style="width:32px; height:32px; border-radius:50%%; flex-shrink:0;">`,
			esc(ev.Author.AvatarURL))
	} else {
		bgColor := timelineIconBg(ev.IconColor)
		fmt.Fprintf(b, `<div style="width:32px; height:32px; border-radius:50%%; background:%s; display:flex; align-items:center; justify-content:center; flex-shrink:0;">`, bgColor)
		fmt.Fprintf(b, `<i class="fa %s" style="color:#fff; font-size:14px;"></i>`, esc(ev.IconClass))
		b.WriteString(`</div>`)
	}

	// Content
	if isMajor {
		b.WriteString(`<div style="flex:1;min-width:0;">`)
	} else {
		b.WriteString(`<div style="min-width:0;">`)
	}
	fmt.Fprintf(b, `<div class="phui-timeline-title" style="font-size:13px; display:flex; align-items:baseline; gap:8px;"><span><strong>%s</strong> %s</span><span style="font-size:12px; color:#6b748c; margin-left:auto; white-space:nowrap;">%s</span></div>`,
		esc(ev.Author.Login), esc(ev.Action), esc(timeAgo(ev.CreatedAt)))
	if isMajor {
		b.WriteString(`<div style="background:#f6f8fa; border:1px solid #e3e4e8; border-radius:4px; padding:12px; font-size:13px; line-height:1.5; overflow-wrap:break-word; word-break:break-word; overflow-x:auto; max-width:100%; margin-top:8px;">`)
		b.WriteString(remarkup.Render(ev.Body))
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)

	b.WriteString(`</div>`) // phui-timeline-event-view
}

func timelineIconBg(color string) string {
	switch color {
	case "green":
		return "#139543"
	case "red":
		return "#c0392b"
	case "blue":
		return "#136cb2"
	case "violet":
		return "#6e5494"
	default:
		return "#6b748c"
	}
}

func renderReviewForm(b *strings.Builder, owner, repo string, number int, sess *auth.Session) {
	b.WriteString(`<div class="mood-review-form" style="margin-top:16px">`)
	fmt.Fprintf(b, `<form method="POST" action="/api/review" data-sigil="workflow">`)
	fmt.Fprintf(b, `<input type="hidden" name="owner" value="%s">`, template.HTMLEscapeString(owner))
	fmt.Fprintf(b, `<input type="hidden" name="repo" value="%s">`, template.HTMLEscapeString(repo))
	fmt.Fprintf(b, `<input type="hidden" name="number" value="%d">`, number)

	b.WriteString(`<textarea name="body" placeholder="Leave a review comment..."></textarea>`)

	b.WriteString(`<div class="form-footer">`)
	b.WriteString(`<span class="pending" id="pending-comment-count"></span>`)
	b.WriteString(`<button type="submit" name="action" value="COMMENT" class="mood-btn mood-btn-default">Comment</button>`)
	b.WriteString(`<button type="submit" name="action" value="APPROVE" class="mood-btn mood-btn-green">Accept</button>`)
	b.WriteString(`<button type="submit" name="action" value="REQUEST_CHANGES" class="mood-btn mood-btn-red">Request Changes</button>`)
	b.WriteString(`</div>`)

	b.WriteString(`</form>`)
	b.WriteString(`</div>`)
}

func renderCurtainReviewersLabels(b *strings.Builder, pr *ghapi.PullRequest, reviews []ghapi.Review) {
	esc := template.HTMLEscapeString

	// Reviewers
	b.WriteString(`<div class="mood-curtain-box">`)
	b.WriteString(`<div class="mood-curtain-title">Reviewers</div>`)
	if len(pr.Reviewers) > 0 {
		for _, u := range pr.Reviewers {
			reviewState := reviewStateForUser(reviews, u.Login)
			b.WriteString(`<div style="display:flex;align-items:center;gap:8px;margin-bottom:8px">`)
			if u.AvatarURL != "" {
				fmt.Fprintf(b, `<img src="%s" style="width:24px;height:24px;border-radius:3px">`, esc(u.AvatarURL))
			}
			fmt.Fprintf(b, `<span style="font-size:13px">%s</span>`, esc(u.Login))
			displayText, displayIcon := reviewStateDisplay(reviewState)
			shadeColor := reviewStateColor(reviewState)
			iconHTML := ""
			if displayIcon != "" {
				iconHTML = fmt.Sprintf(`<span class="phui-icon-view phui-font-fa %s mrs"></span>`, displayIcon)
			}
			fmt.Fprintf(b, `<span class="phui-tag-view phui-tag-shade-%s phui-tag-type-shade" style="margin-left:auto"><span class="phui-tag-core">%s%s</span></span>`,
				shadeColor, iconHTML, esc(displayText))
			b.WriteString(`</div>`)
		}
	} else {
		b.WriteString(`<div style="font-size:13px;color:#6b748c">None assigned</div>`)
	}
	b.WriteString(`</div>`)

	// Labels
	if len(pr.Labels) > 0 {
		b.WriteString(`<div class="mood-curtain-box">`)
		b.WriteString(`<div class="mood-curtain-title">Labels</div>`)
		for _, l := range pr.Labels {
			shade := labelShadeFromColor(l.Color)
			fmt.Fprintf(b, `<span class="phui-tag-view phui-tag-shade-%s phui-tag-type-shade"><span class="phui-tag-core">%s</span></span> `,
				shade, esc(l.Name))
		}
		b.WriteString(`</div>`)
	}
}

func renderCurtainPropertiesActions(b *strings.Builder, pr *ghapi.PullRequest, owner, repo string, number int) {
	esc := template.HTMLEscapeString

	// Properties (Author, Base, Head, Changes)
	b.WriteString(`<div class="mood-curtain-box">`)
	b.WriteString(`<div class="mood-curtain-title">Properties</div>`)
	b.WriteString(`<div class="mood-curtain-prop"><span class="mood-curtain-key">Author</span><span class="mood-curtain-val">`)
	fmt.Fprintf(b, `%s`, esc(pr.Author.Login))
	b.WriteString(`</span></div>`)
	b.WriteString(`<div class="mood-curtain-prop"><span class="mood-curtain-key">Base</span><span class="mood-curtain-val">`)
	fmt.Fprintf(b, `%s`, esc(pr.Base.Ref))
	b.WriteString(`</span></div>`)
	b.WriteString(`<div class="mood-curtain-prop"><span class="mood-curtain-key">Head</span><span class="mood-curtain-val">`)
	fmt.Fprintf(b, `%s`, esc(pr.Head.Ref))
	b.WriteString(`</span></div>`)
	b.WriteString(`<div class="mood-curtain-prop"><span class="mood-curtain-key">Changes</span><span class="mood-curtain-val">`)
	fmt.Fprintf(b, `<span style="color:#139543">+%d</span> / <span style="color:#c0392b">-%d</span> in %d files`, pr.Additions, pr.Deletions, pr.ChangedFiles)
	b.WriteString(`</span></div>`)
	b.WriteString(`</div>`)

	// Actions (only for non-merged PRs)
	if !pr.Merged {
		b.WriteString(`<div class="mood-curtain-actions">`)
		if pr.State != "closed" {
			fmt.Fprintf(b, `<form method="POST" action="/api/merge" style="display:flex;gap:8px;align-items:center;margin-bottom:8px">`)
			fmt.Fprintf(b, `<input type="hidden" name="owner" value="%s">`, esc(owner))
			fmt.Fprintf(b, `<input type="hidden" name="repo" value="%s">`, esc(repo))
			fmt.Fprintf(b, `<input type="hidden" name="number" value="%d">`, number)
			b.WriteString(`<select name="merge_method" style="font-size:12px;padding:4px 6px;border:1px solid #c7ccd9;border-radius:3px;background:#fff;flex:1">`)
			b.WriteString(`<option value="merge">Merge</option><option value="squash">Squash</option><option value="rebase">Rebase</option>`)
			b.WriteString(`</select>`)
			b.WriteString(`<button type="submit" class="mood-btn mood-btn-green" style="flex:1;text-align:center"><span class="phui-icon-view phui-font-fa fa-check-circle mrs"></span>Land Revision</button>`)
			b.WriteString(`</form>`)

			fmt.Fprintf(b, `<form method="POST" action="/api/close">`)
			fmt.Fprintf(b, `<input type="hidden" name="owner" value="%s">`, esc(owner))
			fmt.Fprintf(b, `<input type="hidden" name="repo" value="%s">`, esc(repo))
			fmt.Fprintf(b, `<input type="hidden" name="number" value="%d">`, number)
			fmt.Fprintf(b, `<input type="hidden" name="state" value="closed">`)
			b.WriteString(`<button type="submit" class="mood-btn mood-btn-default" style="width:100%;text-align:center"><span class="phui-icon-view phui-font-fa fa-times-circle mrs"></span>Close</button>`)
			b.WriteString(`</form>`)
		} else {
			fmt.Fprintf(b, `<form method="POST" action="/api/close">`)
			fmt.Fprintf(b, `<input type="hidden" name="owner" value="%s">`, esc(owner))
			fmt.Fprintf(b, `<input type="hidden" name="repo" value="%s">`, esc(repo))
			fmt.Fprintf(b, `<input type="hidden" name="number" value="%d">`, number)
			fmt.Fprintf(b, `<input type="hidden" name="state" value="open">`)
			b.WriteString(`<button type="submit" class="mood-btn mood-btn-green" style="width:100%;text-align:center"><span class="phui-icon-view phui-font-fa fa-refresh mrs"></span>Reopen</button>`)
			b.WriteString(`</form>`)
		}
		b.WriteString(`</div>`)
	}
}

func prStatusBadge(pr *ghapi.PullRequest) (text, color string) {
	if pr.Merged {
		return "Merged", "violet"
	}
	if pr.State == "closed" {
		return "Closed", "red"
	}
	if pr.Draft {
		return "Draft", "grey"
	}
	return "Open", "green"
}

func reviewStateForUser(reviews []ghapi.Review, login string) string {
	// Find the most recent review by this user.
	var latest *ghapi.Review
	for i := range reviews {
		if reviews[i].Author.Login == login {
			if latest == nil || reviews[i].CreatedAt.After(latest.CreatedAt) {
				latest = &reviews[i]
			}
		}
	}
	if latest == nil {
		return ""
	}
	return latest.State
}

func reviewStateDisplay(state string) (text, icon string) {
	switch state {
	case "APPROVED":
		return "Accepted", "fa-check"
	case "CHANGES_REQUESTED":
		return "Changes Requested", "fa-times"
	case "COMMENTED":
		return "Commented", "fa-comment"
	case "DISMISSED":
		return "Dismissed", "fa-ban"
	default:
		return "Waiting", "fa-clock-o"
	}
}

func reviewStateColor(state string) string {
	switch state {
	case "APPROVED":
		return "green"
	case "CHANGES_REQUESTED":
		return "red"
	case "COMMENTED":
		return "blue"
	case "DISMISSED":
		return "grey"
	case "":
		return "orange"
	default:
		return "grey"
	}
}

func labelShadeFromColor(hex string) string {
	if len(hex) < 6 {
		return "blue"
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	bl, _ := strconv.ParseUint(hex[4:6], 16, 8)
	if r > 180 && g < 120 && bl < 120 {
		return "red"
	}
	if g > 180 && r < 120 {
		return "green"
	}
	if bl > 180 && r < 150 {
		return "blue"
	}
	if r > 180 && g > 180 && bl < 120 {
		return "yellow"
	}
	if r > 150 && g > 80 && g < 170 && bl < 100 {
		return "orange"
	}
	if r > 80 && bl > 150 && g < 120 {
		return "violet"
	}
	if r > 200 && g > 200 && bl > 200 {
		return "grey"
	}
	return "blue"
}
