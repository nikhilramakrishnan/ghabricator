package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/templates"

	gh "github.com/google/go-github/v68/github"
)

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")
	if searchType == "" {
		searchType = "pr"
	}

	var buf bytes.Buffer

	// Search form
	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box mlt mlr">`)
	buf.WriteString(`<form method="GET" action="/search" class="phui-form-view" style="padding:12px">`)
	buf.WriteString(`<div style="display:flex;gap:8px;align-items:center">`)
	buf.WriteString(fmt.Sprintf(`<input type="text" name="q" value="%s" class="aphront-form-input" style="flex:1;padding:8px 12px" placeholder="Search pull requests, code, and repositories...">`, template.HTMLEscapeString(query)))
	buf.WriteString(`<select name="type" class="aphront-form-input" style="padding:8px">`)
	for _, opt := range []struct{ val, label string }{
		{"pr", "Pull Requests"},
		{"code", "Code"},
		{"repo", "Repositories"},
	} {
		sel := ""
		if searchType == opt.val {
			sel = " selected"
		}
		buf.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, opt.val, sel, opt.label))
	}
	buf.WriteString(`</select>`)
	buf.WriteString(`<button type="submit" class="phui-button-view button-green" style="padding:8px 16px">`)
	buf.WriteString(`<span class="phui-icon-view phui-font-fa fa-search"></span> Search`)
	buf.WriteString(`</button>`)
	buf.WriteString(`</div>`)
	buf.WriteString(`</form>`)
	buf.WriteString(`</div>`)

	// Only run search if there's a query
	if query != "" {
		switch searchType {
		case "code":
			buf.WriteString(s.searchCode(r, client, query))
		case "repo":
			buf.WriteString(s.searchRepos(r, client, query))
		default:
			buf.WriteString(s.searchPRsForSearch(r, client, query))
		}
	}

	templates.RenderPage(w, templates.PageData{
		Title:         "Search",
		Theme:         theme,
		HeaderTitle:   template.HTML("Search"),
		HeaderIcon:    "fa-search",
		Content:       template.HTML(buf.String()),
		NavActive:     "search",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Search"},
		},
	})
}

func (s *Server) searchPRsForSearch(r *http.Request, client *gh.Client, query string) string {
	prs := s.searchPRs(r, client, "is:pr "+query)
	return renderPRSection("Pull Request Results", "fa-code-fork", prs)
}

func (s *Server) searchCode(r *http.Request, client *gh.Client, query string) string {
	result, _, err := client.Search.Code(r.Context(), query, &gh.SearchOptions{
		ListOptions: gh.ListOptions{PerPage: 25},
	})
	if err != nil {
		log.Printf("search code error: %v", err)
		return renderSearchError("Code search failed.")
	}

	var buf bytes.Buffer
	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box mlt mlr">`)
	buf.WriteString(`<div class="phui-header-shell"><div class="phui-header-view"><h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-file-code-o"></span>Code Results</h1></div></div>`)

	if len(result.CodeResults) == 0 {
		buf.WriteString(`<div class="phui-info-view phui-info-severity-nodata"><span class="phui-icon-view phui-font-fa fa-search mrs"></span>No code results found.</div>`)
		buf.WriteString(`</div>`)
		return buf.String()
	}

	buf.WriteString(`<ul class="phui-oi-list-view">`)
	for _, cr := range result.CodeResults {
		repo := cr.GetRepository().GetFullName()
		path := cr.GetPath()
		href := fmt.Sprintf("/pr/%s", repo) // link to repo view

		buf.WriteString(`<li class="phui-oi phui-oi-bar-color-violet">`)
		buf.WriteString(`<div class="phui-oi-frame"><div class="phui-oi-frame-content">`)
		buf.WriteString(`<div class="phui-oi-image-icon"><span class="phui-icon-view phui-font-fa fa-file-code-o"></span></div>`)
		buf.WriteString(`<div class="phui-oi-content-box"><div class="phui-oi-table"><div class="phui-oi-table-row"><div class="phui-oi-col1">`)

		// Name: owner/repo/path
		buf.WriteString(`<div class="phui-oi-name">`)
		buf.WriteString(fmt.Sprintf(`<a href="%s" class="phui-oi-link">%s/%s</a>`,
			template.HTMLEscapeString(href),
			template.HTMLEscapeString(repo),
			template.HTMLEscapeString(path)))
		buf.WriteString(`</div>`)

		// Code fragment
		buf.WriteString(`<div class="phui-oi-content">`)
		if len(cr.TextMatches) > 0 {
			for _, tm := range cr.TextMatches {
				fragment := tm.GetFragment()
				if fragment != "" {
					buf.WriteString(`<pre class="PhabricatorMonospaced" style="font-size:12px;padding:4px 8px;background:rgba(55,55,55,.04);border-radius:3px;white-space:pre-wrap;word-break:break-all;margin:4px 0">`)
					buf.WriteString(template.HTMLEscapeString(fragment))
					buf.WriteString(`</pre>`)
				}
			}
		}
		buf.WriteString(`</div>`)

		buf.WriteString(`</div></div></div>`) // col1, table-row, table
		buf.WriteString(`</div>`)             // content-box
		buf.WriteString(`</div></div>`)       // frame-content, frame
		buf.WriteString(`</li>`)
	}
	buf.WriteString(`</ul>`)
	buf.WriteString(`</div>`)
	return buf.String()
}

func (s *Server) searchRepos(r *http.Request, client *gh.Client, query string) string {
	result, _, err := client.Search.Repositories(r.Context(), query, &gh.SearchOptions{
		Sort:        "stars",
		Order:       "desc",
		ListOptions: gh.ListOptions{PerPage: 25},
	})
	if err != nil {
		log.Printf("search repos error: %v", err)
		return renderSearchError("Repository search failed.")
	}

	var buf bytes.Buffer
	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box mlt mlr">`)
	buf.WriteString(`<div class="phui-header-shell"><div class="phui-header-view"><h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-database"></span>Repository Results</h1></div></div>`)

	if len(result.Repositories) == 0 {
		buf.WriteString(`<div class="phui-info-view phui-info-severity-nodata"><span class="phui-icon-view phui-font-fa fa-search mrs"></span>No repositories found.</div>`)
		buf.WriteString(`</div>`)
		return buf.String()
	}

	buf.WriteString(`<ul class="phui-oi-list-view">`)
	for _, repo := range result.Repositories {
		fullName := repo.GetFullName()
		href := fmt.Sprintf("/pr/%s", fullName)
		desc := repo.GetDescription()
		stars := repo.GetStargazersCount()
		lang := repo.GetLanguage()

		avatarURL := ""
		if repo.Owner != nil {
			avatarURL = repo.Owner.GetAvatarURL()
		}

		imageClass := "phui-oi-with-image-icon"
		if avatarURL != "" {
			imageClass = "phui-oi-with-image"
		}

		buf.WriteString(fmt.Sprintf(`<li class="phui-oi phui-oi-bar-color-blue %s">`, imageClass))
		buf.WriteString(`<div class="phui-oi-frame"><div class="phui-oi-frame-content">`)

		if avatarURL != "" {
			buf.WriteString(fmt.Sprintf(`<div class="phui-oi-image" style="background-image:url(%s)"></div>`, template.HTMLEscapeString(avatarURL)))
		} else {
			buf.WriteString(`<div class="phui-oi-image-icon"><span class="phui-icon-view phui-font-fa fa-database"></span></div>`)
		}

		buf.WriteString(`<div class="phui-oi-content-box"><div class="phui-oi-table"><div class="phui-oi-table-row"><div class="phui-oi-col1">`)

		// Name
		buf.WriteString(`<div class="phui-oi-name">`)
		buf.WriteString(fmt.Sprintf(`<a href="%s" class="phui-oi-link">%s</a>`,
			template.HTMLEscapeString(href),
			template.HTMLEscapeString(fullName)))
		buf.WriteString(`</div>`)

		// Attributes
		buf.WriteString(`<div class="phui-oi-content"><ul class="phui-oi-attributes">`)
		if desc != "" {
			buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute">%s</li>`, template.HTMLEscapeString(desc)))
		}
		buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-star mrs"></span>%d</li>`, stars))
		if lang != "" {
			buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute">%s</li>`, template.HTMLEscapeString(lang)))
		}
		buf.WriteString(`</ul></div>`)

		buf.WriteString(`</div></div></div>`) // col1, table-row, table
		buf.WriteString(`</div>`)             // content-box
		buf.WriteString(`</div></div>`)       // frame-content, frame
		buf.WriteString(`</li>`)
	}
	buf.WriteString(`</ul>`)
	buf.WriteString(`</div>`)
	return buf.String()
}

func renderSearchError(msg string) string {
	return fmt.Sprintf(`<div class="phui-box phui-box-border phui-object-box mlt mlr"><div class="phui-info-view phui-info-severity-warning"><span class="phui-icon-view phui-font-fa fa-exclamation-triangle mrs"></span>%s</div></div>`, template.HTMLEscapeString(msg))
}
