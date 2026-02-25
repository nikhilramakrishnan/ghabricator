package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
	"github.com/nikhilr/ghabricator/internal/templates"
)

func (s *Server) handlePasteList(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	gists, err := ghapi.ListGists(r.Context(), client)
	if err != nil {
		log.Printf("paste list error: %v", err)
		gists = nil
	}

	var buf bytes.Buffer

	// Paste list
	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box">`)
	buf.WriteString(`<div class="phui-header-shell" style="display:flex;align-items:center;justify-content:space-between">`)
	buf.WriteString(`<div class="phui-header-view"><h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-clipboard"></span>Recent Pastes</h1></div>`)
	buf.WriteString(`<a href="/paste/new" class="button-green phui-button-view" style="margin-right:8px"><span class="phui-icon-view phui-font-fa fa-plus mrs"></span>Create Paste</a>`)
	buf.WriteString(`</div>`)

	if len(gists) == 0 {
		buf.WriteString(`<div class="phui-info-view phui-info-severity-nodata">No pastes found. Create one!</div>`)
	} else {
		buf.WriteString(`<ul class="phui-oi-list-view">`)
		for _, g := range gists {
			title := g.Description
			if title == "" {
				title = "(untitled)"
			}

			// Get first file info
			lang := ""
			for _, f := range g.Files {
				if f.Language != "" {
					lang = f.Language
				}
				break
			}

			visIcon := "fa-globe"
			visLabel := "Public"
			if !g.Public {
				visIcon = "fa-lock"
				visLabel = "Secret"
			}

			buf.WriteString(`<li class="phui-oi phui-oi-standard">`)
			buf.WriteString(`<div class="phui-oi-frame">`)

			// Icon
			buf.WriteString(`<div class="phui-oi-image-icon">`)
			buf.WriteString(`<span class="phui-icon-view phui-font-fa fa-file-code-o"></span>`)
			buf.WriteString(`</div>`)

			// Content
			buf.WriteString(`<div class="phui-oi-content-box">`)
			buf.WriteString(`<div class="phui-oi-name">`)
			fmt.Fprintf(&buf, `<a href="/paste/%s">%s</a>`, template.HTMLEscapeString(g.ID), template.HTMLEscapeString(title))
			buf.WriteString(`</div>`)

			// Attributes
			buf.WriteString(`<div class="phui-oi-attributes">`)
			if lang != "" {
				fmt.Fprintf(&buf, `<span class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-code"></span> %s</span>`, template.HTMLEscapeString(lang))
			}
			fmt.Fprintf(&buf, `<span class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa %s"></span> %s</span>`, visIcon, visLabel)
			fmt.Fprintf(&buf, `<span class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-clock-o"></span> %s</span>`, timeAgo(g.CreatedAt))
			buf.WriteString(`</div>`)

			buf.WriteString(`</div>`) // content-box
			buf.WriteString(`</div>`) // frame
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	}
	buf.WriteString(`</div>`) // object-box

	templates.RenderPage(w, templates.PageData{
		Title:         "Paste",
		Theme:         theme,
		HeaderTitle:   "Paste",
		HeaderIcon:    "fa-clipboard",
		Content:       template.HTML(buf.String()),
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Paste"},
		},
	})
}

func (s *Server) handlePasteNew(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	var buf bytes.Buffer
	buf.WriteString(`<div class="phui-object-box">`)
	buf.WriteString(`<div class="phui-header-shell"><div class="phui-header-view"><h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-plus"></span>Create Paste</h1></div></div>`)

	buf.WriteString(`<form method="POST" action="/paste/save">`)
	buf.WriteString(`<div class="phui-form-view">`)

	// Title
	buf.WriteString(`<div class="aphront-form-control">`)
	buf.WriteString(`<label class="aphront-form-label">Title</label>`)
	buf.WriteString(`<div class="aphront-form-input"><input type="text" name="title" placeholder="Paste title"></div>`)
	buf.WriteString(`</div>`)

	// Language
	buf.WriteString(`<div class="aphront-form-control">`)
	buf.WriteString(`<label class="aphront-form-label">Language</label>`)
	buf.WriteString(`<div class="aphront-form-input"><select name="language">`)
	for _, lang := range pasteLanguages {
		fmt.Fprintf(&buf, `<option value="%s">%s</option>`, lang.ext, template.HTMLEscapeString(lang.name))
	}
	buf.WriteString(`</select></div>`)
	buf.WriteString(`</div>`)

	// Content
	buf.WriteString(`<div class="aphront-form-control aphront-form-control-textarea">`)
	buf.WriteString(`<label class="aphront-form-label">Content</label>`)
	buf.WriteString(`<div class="aphront-form-input"><textarea name="content" rows="20" class="PhabricatorMonospaced aphront-textarea-very-tall" placeholder="Paste content here..." required></textarea></div>`)
	buf.WriteString(`</div>`)

	// Visibility
	buf.WriteString(`<div class="aphront-form-control">`)
	buf.WriteString(`<label class="aphront-form-label">Visibility</label>`)
	buf.WriteString(`<div class="aphront-form-input"><select name="visibility">`)
	buf.WriteString(`<option value="secret">Secret</option>`)
	buf.WriteString(`<option value="public">Public</option>`)
	buf.WriteString(`</select></div>`)
	buf.WriteString(`</div>`)

	// Submit
	buf.WriteString(`<div class="aphront-form-control aphront-form-control-submit">`)
	buf.WriteString(`<div class="aphront-form-input">`)
	buf.WriteString(`<button type="submit" class="phui-button-view button-green">`)
	buf.WriteString(`<span class="phui-icon-view phui-font-fa fa-save"></span> Create Paste</button>`)
	buf.WriteString(`</div></div>`)

	buf.WriteString(`</div>`) // form-view
	buf.WriteString(`</form>`)
	buf.WriteString(`</div>`) // object-box

	templates.RenderPage(w, templates.PageData{
		Title:         "Create Paste",
		Theme:         theme,
		HeaderTitle:   "Create Paste",
		HeaderIcon:    "fa-plus",
		Content:       template.HTML(buf.String()),
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Paste", Href: "/paste"},
			{Name: "Create"},
		},
	})
}

func (s *Server) handlePasteSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderError(w, r, "Invalid Request", "Could not parse form data.", http.StatusBadRequest)
		return
	}

	client := auth.GitHubClientFromContext(r.Context())

	title := r.FormValue("title")
	language := r.FormValue("language")
	content := r.FormValue("content")
	visibility := r.FormValue("visibility")

	if content == "" {
		s.renderError(w, r, "Invalid Request", "Paste content is required.", http.StatusBadRequest)
		return
	}

	if title == "" {
		title = "Untitled Paste"
	}

	filename := "paste"
	if language != "" {
		filename = "paste." + language
	}

	public := visibility == "public"

	gist, err := ghapi.CreateGist(r.Context(), client, title, filename, content, public)
	if err != nil {
		log.Printf("paste save error: %v", err)
		s.renderError(w, r, "Save Error", fmt.Sprintf("Could not create paste: %v", err), http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/paste/"+gist.ID, http.StatusSeeOther)
}

func (s *Server) handlePasteView(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)
	gistID := r.PathValue("id")

	gist, err := ghapi.FetchGist(r.Context(), client, gistID)
	if err != nil {
		s.renderError(w, r, "Not Found", fmt.Sprintf("Paste not found: %v", err), http.StatusNotFound)
		return
	}

	title := gist.Description
	if title == "" {
		title = "(untitled)"
	}

	var buf bytes.Buffer

	// Iterate over files in the gist
	for name, f := range gist.Files {
		buf.WriteString(`<div class="phui-object-box">`)

		// File header
		buf.WriteString(`<div class="phui-header-shell"><div class="phui-header-view">`)
		fmt.Fprintf(&buf, `<h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-file-code-o"></span>%s`, template.HTMLEscapeString(name))
		if f.Language != "" {
			fmt.Fprintf(&buf, ` <span class="phui-tag-view phui-tag-shade-blue phui-tag-type-shade"><span class="phui-tag-core">%s</span></span>`, template.HTMLEscapeString(f.Language))
		}
		buf.WriteString(`</h1></div></div>`)

		// Syntax-highlighted code
		content := f.Content
		lines := strings.Split(content, "\n")
		highlighted := diff.HighlightLines(name, lines)

		buf.WriteString(`<div class="phabricator-source-code-container">`)
		buf.WriteString(`<table class="phabricator-source-code-view remarkup-code PhabricatorMonospaced chroma">`)
		for i, hl := range highlighted {
			lineNum := i + 1
			buf.WriteString(`<tr>`)
			fmt.Fprintf(&buf, `<th class="phabricator-source-line"><span>%d</span></th>`, lineNum)
			fmt.Fprintf(&buf, `<td class="phabricator-source-code">%s</td>`, hl)
			buf.WriteString(`</tr>`)
		}
		buf.WriteString(`</table>`)
		buf.WriteString(`</div>`)

		buf.WriteString(`</div>`) // object-box
	}

	// Curtain with metadata
	var curtain strings.Builder
	curtain.WriteString(`<div class="phui-box phui-box-border phui-object-box phui-curtain-view">`)

	// Author
	if gist.Owner.Login != "" {
		curtain.WriteString(`<div class="phui-curtain-panel">`)
		curtain.WriteString(`<div class="phui-curtain-panel-header">Author</div>`)
		curtain.WriteString(`<div class="phui-curtain-panel-body">`)
		if gist.Owner.AvatarURL != "" {
			fmt.Fprintf(&curtain, `<img src="%s" class="phui-head-small" alt="">`,
				template.HTMLEscapeString(gist.Owner.AvatarURL))
		}
		fmt.Fprintf(&curtain, ` %s`, template.HTMLEscapeString(gist.Owner.Login))
		curtain.WriteString(`</div></div>`)
	}

	// Details
	curtain.WriteString(`<div class="phui-curtain-panel">`)
	curtain.WriteString(`<div class="phui-curtain-panel-header">Details</div>`)
	curtain.WriteString(`<div class="phui-curtain-panel-body">`)

	visIcon := "fa-globe"
	visLabel := "Public"
	if !gist.Public {
		visIcon = "fa-lock"
		visLabel = "Secret"
	}
	fmt.Fprintf(&curtain, `<div><span class="phui-icon-view phui-font-fa %s"></span> %s</div>`, visIcon, visLabel)
	fmt.Fprintf(&curtain, `<div>Created %s</div>`, timeAgo(gist.CreatedAt))
	if !gist.UpdatedAt.Equal(gist.CreatedAt) {
		fmt.Fprintf(&curtain, `<div>Updated %s</div>`, timeAgo(gist.UpdatedAt))
	}
	curtain.WriteString(`</div></div>`)

	// Actions
	curtain.WriteString(`<div class="phui-curtain-panel">`)
	curtain.WriteString(`<div class="phui-curtain-panel-header">Actions</div>`)
	curtain.WriteString(`<div class="phui-curtain-panel-body">`)
	curtain.WriteString(`<ul class="phabricator-action-list-view">`)
	fmt.Fprintf(&curtain, `<li class="phabricator-action-view action-has-icon"><a href="%s" target="_blank" class="phabricator-action-view-item"><span class="phabricator-action-view-icon phui-icon-view phui-font-fa fa-github"></span>View on GitHub</a></li>`,
		template.HTMLEscapeString(gist.HTMLURL))
	curtain.WriteString(`</ul>`)
	curtain.WriteString(`</div></div>`)

	curtain.WriteString(`</div>`)

	templates.RenderPage(w, templates.PageData{
		Title:         fmt.Sprintf("P%s: %s", gistID[:8], title),
		Theme:         theme,
		HeaderTitle:   template.HTML(template.HTMLEscapeString(title)),
		HeaderIcon:    "fa-clipboard",
		Content:       template.HTML(buf.String()),
		Curtain:       template.HTML(curtain.String()),
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Paste", Href: "/paste"},
			{Name: fmt.Sprintf("P%s", gistID[:8])},
		},
	})
}

var pasteLanguages = []struct {
	name string
	ext  string
}{
	{"Plain Text", "txt"},
	{"Go", "go"},
	{"Python", "py"},
	{"JavaScript", "js"},
	{"TypeScript", "ts"},
	{"Rust", "rs"},
	{"C", "c"},
	{"C++", "cpp"},
	{"Java", "java"},
	{"Ruby", "rb"},
	{"PHP", "php"},
	{"Shell", "sh"},
	{"SQL", "sql"},
	{"HTML", "html"},
	{"CSS", "css"},
	{"JSON", "json"},
	{"YAML", "yaml"},
	{"Markdown", "md"},
	{"XML", "xml"},
	{"Diff", "diff"},
}
