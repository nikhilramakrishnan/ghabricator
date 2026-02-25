package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/templates"

	gh "github.com/google/go-github/v68/github"
)

type dashboardPR struct {
	Title     string
	Repo      string
	Number    int
	URL       string
	Author    string
	AvatarURL string
	UpdatedAt string
	Draft     bool
	Labels    []string
	Reviewers []dashboardUser
}

type dashboardUser struct {
	Login     string
	AvatarURL string
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)
	login := sess.Login

	// Search for PRs authored by user and PRs requesting user's review
	authored := s.searchPRs(r, client, fmt.Sprintf("is:open is:pr author:%s", login))
	reviewing := s.searchPRs(r, client, fmt.Sprintf("is:open is:pr review-requested:%s", login))

	var buf bytes.Buffer
	buf.WriteString(renderPRSection("Authored", "fa-pencil", authored))
	buf.WriteString(renderPRSection("Review Requested", "fa-eye", reviewing))

	templates.RenderPage(w, templates.PageData{
		Title:         "Dashboard",
		Theme:         theme,
		HeaderTitle:   template.HTML("Active Revisions"),
		HeaderIcon:    "fa-th-list",
		Content:       template.HTML(buf.String()),
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Dashboard"},
		},
	})
}

func (s *Server) searchPRs(r *http.Request, client *gh.Client, query string) []dashboardPR {
	result, _, err := client.Search.Issues(r.Context(), query, &gh.SearchOptions{
		Sort:  "updated",
		Order: "desc",
		ListOptions: gh.ListOptions{PerPage: 25},
	})
	if err != nil {
		log.Printf("dashboard search error: %v", err)
		return nil
	}

	var prs []dashboardPR
	for _, issue := range result.Issues {
		if issue.PullRequestLinks == nil {
			continue
		}
		repo := ""
		if issue.Repository != nil {
			repo = issue.Repository.GetFullName()
		} else if issue.RepositoryURL != nil {
			// Parse owner/repo from URL like https://api.github.com/repos/owner/repo
			repo = extractRepoFromURL(issue.GetRepositoryURL())
		}

		pr := dashboardPR{
			Title:     issue.GetTitle(),
			Repo:      repo,
			Number:    issue.GetNumber(),
			Author:    issue.GetUser().GetLogin(),
			AvatarURL: issue.GetUser().GetAvatarURL(),
			UpdatedAt: timeAgo(issue.GetUpdatedAt().Time),
			Draft:     issue.IsPullRequest() && issue.GetDraft(),
		}

		// Build link
		if repo != "" {
			pr.URL = fmt.Sprintf("/pr/%s/%d", repo, issue.GetNumber())
		}

		for _, l := range issue.Labels {
			pr.Labels = append(pr.Labels, l.GetName())
		}
		for _, a := range issue.Assignees {
			pr.Reviewers = append(pr.Reviewers, dashboardUser{
				Login:     a.GetLogin(),
				AvatarURL: a.GetAvatarURL(),
			})
		}

		prs = append(prs, pr)
	}
	return prs
}

func renderPRSection(title, icon string, prs []dashboardPR) string {
	var buf bytes.Buffer
	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box mlt mlr">`)
	buf.WriteString(fmt.Sprintf(`<div class="phui-header-shell"><div class="phui-header-view"><h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa %s"></span>%s</h1></div></div>`, icon, title))

	if len(prs) == 0 {
		buf.WriteString(`<div class="phui-info-view phui-info-severity-nodata"><span class="phui-icon-view phui-font-fa fa-inbox mrs"></span>No open pull requests.</div>`)
		buf.WriteString(`</div>`)
		return buf.String()
	}

	buf.WriteString(`<ul class="phui-oi-list-view">`)
	for _, pr := range prs {
		imageClass := "phui-oi-with-image-icon"
		if pr.AvatarURL != "" {
			imageClass = "phui-oi-with-image"
		}

		href := "#"
		if pr.URL != "" {
			href = pr.URL
		}

		barClass := "phui-oi-bar-color-blue"
		if pr.Draft {
			barClass = "phui-oi-bar-color-grey"
		}
		buf.WriteString(fmt.Sprintf(`<li class="phui-oi %s %s">`, barClass, imageClass))
		buf.WriteString(`<div class="phui-oi-frame">`)
		buf.WriteString(`<div class="phui-oi-frame-content">`)

		// Image
		if pr.AvatarURL != "" {
			buf.WriteString(fmt.Sprintf(`<div class="phui-oi-image" style="background-image:url(%s)"></div>`, template.HTMLEscapeString(pr.AvatarURL)))
		} else {
			buf.WriteString(`<div class="phui-oi-image-icon"><span class="phui-icon-view phui-font-fa fa-code-fork"></span></div>`)
		}

		// Content box with table layout
		buf.WriteString(`<div class="phui-oi-content-box">`)
		buf.WriteString(`<div class="phui-oi-table"><div class="phui-oi-table-row"><div class="phui-oi-col1">`)

		// Name
		buf.WriteString(`<div class="phui-oi-name">`)
		buf.WriteString(fmt.Sprintf(`<a href="%s" class="phui-oi-link">%s</a>`, template.HTMLEscapeString(href), template.HTMLEscapeString(pr.Title)))
		if pr.Draft {
			buf.WriteString(` <span class="phui-tag-view phui-tag-shade-grey phui-tag-type-shade"><span class="phui-tag-core">Draft</span></span>`)
		}
		buf.WriteString(`</div>`)

		// Attributes
		buf.WriteString(`<div class="phui-oi-content">`)
		buf.WriteString(`<ul class="phui-oi-attributes">`)
		if pr.Author != "" {
			buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-user mrs"></span>%s</li>`, template.HTMLEscapeString(pr.Author)))
		}
		if pr.Repo != "" {
			buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-github"></span> %s#%d</li>`, template.HTMLEscapeString(pr.Repo), pr.Number))
		}
		buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-clock-o"></span> %s</li>`, pr.UpdatedAt))
		if len(pr.Labels) > 0 {
			buf.WriteString(`<li class="phui-oi-attribute">`)
			for _, label := range pr.Labels {
				buf.WriteString(fmt.Sprintf(`<span class="phui-tag-view phui-tag-shade-blue phui-tag-type-shade"><span class="phui-tag-core">%s</span></span>`, template.HTMLEscapeString(label)))
			}
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
		buf.WriteString(`</div>`) // content

		buf.WriteString(`</div></div></div>`) // col1, table-row, table
		buf.WriteString(`</div>`) // content-box

		buf.WriteString(`</div>`) // frame-content

		// Reviewer handle icons
		if len(pr.Reviewers) > 0 {
			buf.WriteString(`<div class="phui-oi-handle-icons">`)
			for _, rev := range pr.Reviewers {
				if rev.AvatarURL != "" {
					buf.WriteString(fmt.Sprintf(`<span class="phui-oi-handle-icon" title="%s" style="background-image:url(%s)"></span>`, template.HTMLEscapeString(rev.Login), template.HTMLEscapeString(rev.AvatarURL)))
				}
			}
			buf.WriteString(`</div>`)
		}

		buf.WriteString(`</div>`) // frame
		buf.WriteString(`</li>`)
	}
	buf.WriteString(`</ul>`)
	buf.WriteString(`</div>`) // object-box
	return buf.String()
}

func extractRepoFromURL(url string) string {
	// https://api.github.com/repos/owner/repo → owner/repo
	const prefix = "/repos/"
	idx := 0
	for i := 0; i <= len(url)-len(prefix); i++ {
		if url[i:i+len(prefix)] == prefix {
			idx = i + len(prefix)
			break
		}
	}
	if idx > 0 {
		return url[idx:]
	}
	return ""
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(math.Round(d.Minutes()))
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(math.Round(d.Hours()))
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		days := int(math.Round(d.Hours() / 24))
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
