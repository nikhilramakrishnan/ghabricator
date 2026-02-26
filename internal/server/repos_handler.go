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

func (s *Server) handleRepoList(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	// Fetch user's repos
	repos, _, err := client.Repositories.List(r.Context(), "", &gh.RepositoryListOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: gh.ListOptions{PerPage: 50},
	})
	if err != nil {
		log.Printf("repo list error: %v", err)
		repos = nil
	}

	var buf bytes.Buffer

	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box">`)
	buf.WriteString(`<div class="phui-header-shell"><div class="phui-header-view"><h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-database"></span>Your Repositories</h1></div></div>`)

	if len(repos) == 0 {
		buf.WriteString(`<div class="phui-info-view phui-info-severity-nodata"><span class="phui-icon-view phui-font-fa fa-inbox mrs"></span>No repositories found.</div>`)
	} else {
		buf.WriteString(`<ul class="phui-oi-list-view">`)
		for _, repo := range repos {
			fullName := repo.GetFullName()
			href := "/repo/" + fullName
			desc := repo.GetDescription()
			lang := repo.GetLanguage()
			stars := repo.GetStargazersCount()
			forks := repo.GetForksCount()

			barColor := "blue"
			if repo.GetFork() {
				barColor = "violet"
			}
			if repo.GetArchived() {
				barColor = "grey"
			}

			avatarURL := ""
			if repo.Owner != nil {
				avatarURL = repo.Owner.GetAvatarURL()
			}

			imageClass := "phui-oi-with-image-icon"
			if avatarURL != "" {
				imageClass = "phui-oi-with-image"
			}

			buf.WriteString(fmt.Sprintf(`<li class="phui-oi phui-oi-bar-color-%s %s">`, barColor, imageClass))
			buf.WriteString(`<div class="phui-oi-frame"><div class="phui-oi-frame-content">`)

			if avatarURL != "" {
				buf.WriteString(fmt.Sprintf(`<div class="phui-oi-image" style="background-image:url(%s)"></div>`, template.HTMLEscapeString(avatarURL)))
			} else {
				buf.WriteString(`<div class="phui-oi-image-icon"><span class="phui-icon-view phui-font-fa fa-database"></span></div>`)
			}

			buf.WriteString(`<div class="phui-oi-content-box"><div class="phui-oi-table"><div class="phui-oi-table-row"><div class="phui-oi-col1">`)

			// Name + tags
			buf.WriteString(`<div class="phui-oi-name">`)
			buf.WriteString(fmt.Sprintf(`<a href="%s" class="phui-oi-link">%s</a>`, template.HTMLEscapeString(href), template.HTMLEscapeString(fullName)))
			if repo.GetPrivate() {
				buf.WriteString(` <span class="phui-tag-view phui-tag-shade-orange phui-tag-type-shade"><span class="phui-tag-core"><span class="phui-icon-view phui-font-fa fa-lock mrs"></span>Private</span></span>`)
			}
			if repo.GetFork() {
				buf.WriteString(` <span class="phui-tag-view phui-tag-shade-violet phui-tag-type-shade"><span class="phui-tag-core">Fork</span></span>`)
			}
			if repo.GetArchived() {
				buf.WriteString(` <span class="phui-tag-view phui-tag-shade-grey phui-tag-type-shade"><span class="phui-tag-core">Archived</span></span>`)
			}
			buf.WriteString(`</div>`)

			// Description + attributes
			buf.WriteString(`<div class="phui-oi-content"><ul class="phui-oi-attributes">`)
			if desc != "" {
				buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute">%s</li>`, template.HTMLEscapeString(desc)))
			}
			if lang != "" {
				buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-circle mrs" style="font-size:10px;color:%s"></span>%s</li>`, langColor(lang), template.HTMLEscapeString(lang)))
			}
			if stars > 0 {
				buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-star mrs"></span>%d</li>`, stars))
			}
			if forks > 0 {
				buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-code-fork mrs"></span>%d</li>`, forks))
			}
			buf.WriteString(fmt.Sprintf(`<li class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-clock-o mrs"></span>%s</li>`, timeAgo(repo.GetUpdatedAt().Time)))
			buf.WriteString(`</ul></div>`)

			buf.WriteString(`</div></div></div>`) // col1, table-row, table
			buf.WriteString(`</div>`)             // content-box
			buf.WriteString(`</div></div>`)       // frame-content, frame
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	}
	buf.WriteString(`</div>`)

	templates.RenderPage(w, templates.PageData{
		Title:         "Repositories",
		Theme:         theme,
		HeaderTitle:   template.HTML("Repositories"),
		HeaderIcon:    "fa-database",
		Content:       template.HTML(buf.String()),
		NavActive:     "repos",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Repositories"},
		},
	})
}

func langColor(lang string) string {
	colors := map[string]string{
		"Go":         "#00ADD8",
		"Python":     "#3572A5",
		"JavaScript": "#f1e05a",
		"TypeScript": "#2b7489",
		"Rust":       "#dea584",
		"Java":       "#b07219",
		"C":          "#555555",
		"C++":        "#f34b7d",
		"Ruby":       "#701516",
		"PHP":        "#4F5D95",
		"Shell":      "#89e051",
		"HTML":       "#e34c26",
		"CSS":        "#563d7c",
		"Swift":      "#ffac45",
		"Kotlin":     "#A97BFF",
		"Dart":       "#00B4AB",
		"Lua":        "#000080",
	}
	if c, ok := colors[lang]; ok {
		return c
	}
	return "#6b748c"
}
