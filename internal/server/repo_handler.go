package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/diff"
	ghapi "github.com/nikhilr/ghabricator/internal/github"
	"github.com/nikhilr/ghabricator/internal/templates"
)

func (s *Server) handleRepoView(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	path := r.PathValue("path")
	ref := r.URL.Query().Get("ref")

	// Fetch repo info for default branch
	repoInfo, err := ghapi.FetchRepoInfo(r.Context(), client, owner, repo)
	if err != nil {
		log.Printf("repo info error: %v", err)
		http.Error(w, fmt.Sprintf("repository not found: %v", err), http.StatusNotFound)
		return
	}
	if ref == "" {
		ref = repoInfo.DefaultBranch
	}

	fullName := template.HTMLEscapeString(owner + "/" + repo)

	// Try fetching as directory first
	entries, dirErr := ghapi.FetchRepoTree(r.Context(), client, owner, repo, ref, path)
	if dirErr == nil {
		s.renderDirView(w, owner, repo, ref, path, fullName, entries, repoInfo, sess, theme)
		return
	}

	// Try as file
	file, fileErr := ghapi.FetchFileContent(r.Context(), client, owner, repo, ref, path)
	if fileErr == nil {
		s.renderFileView(w, owner, repo, ref, path, fullName, file, repoInfo, sess, theme)
		return
	}

	log.Printf("repo view error: dir=%v file=%v", dirErr, fileErr)
	http.Error(w, "path not found", http.StatusNotFound)
}

func (s *Server) renderDirView(w http.ResponseWriter, owner, repo, ref, path, fullName string, entries []ghapi.RepoEntry, repoInfo *ghapi.RepoInfo, sess *auth.Session, theme string) {
	var buf bytes.Buffer

	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box">`)

	buf.WriteString(`<table class="aphront-table-view">`)
	buf.WriteString(`<tr><th>Name</th><th>Type</th><th>Size</th></tr>`)

	// Parent directory link
	if path != "" {
		parentPath := filepath.Dir(path)
		parentHref := fmt.Sprintf("/repo/%s/%s", owner, repo)
		if parentPath != "." {
			parentHref = fmt.Sprintf("/repo/%s/%s/%s", owner, repo, parentPath)
		}
		parentHref += "?ref=" + template.URLQueryEscaper(ref)
		buf.WriteString(`<tr>`)
		fmt.Fprintf(&buf, `<td><span class="phui-icon-view phui-font-fa fa-level-up mrs"></span><a href="%s">..</a></td>`, template.HTMLEscapeString(parentHref))
		buf.WriteString(`<td></td><td></td>`)
		buf.WriteString(`</tr>`)
	}

	// Directories first, then files
	for _, pass := range []string{"dir", "file"} {
		for _, e := range entries {
			if e.Type != pass {
				continue
			}
			href := fmt.Sprintf("/repo/%s/%s/%s?ref=%s", owner, repo, e.Path, template.URLQueryEscaper(ref))
			escapedHref := template.HTMLEscapeString(href)
			escapedName := template.HTMLEscapeString(e.Name)

			buf.WriteString(`<tr>`)
			if e.Type == "dir" {
				fmt.Fprintf(&buf, `<td><span class="phui-icon-view phui-font-fa fa-folder mrs"></span><a href="%s">%s/</a></td>`, escapedHref, escapedName)
				buf.WriteString(`<td>dir</td><td>&mdash;</td>`)
			} else {
				icon := diff.FileIcon(e.Name)
				fmt.Fprintf(&buf, `<td><span class="phui-icon-view phui-font-fa %s mrs"></span><a href="%s">%s</a></td>`, icon, escapedHref, escapedName)
				buf.WriteString(`<td>file</td>`)
				fmt.Fprintf(&buf, `<td>%s</td>`, formatSize(e.Size))
			}
			buf.WriteString(`</tr>`)
		}
	}

	buf.WriteString(`</table>`)
	buf.WriteString(`</div>`)

	// Curtain with repo info
	curtain := buildRepoCurtain(repoInfo)

	dirName := filepath.Base(path)
	if path == "" {
		dirName = owner + "/" + repo
	}

	templates.RenderPage(w, templates.PageData{
		Title:         dirName + " - Diffusion",
		Theme:         theme,
		Content:       template.HTML(buf.String()),
		Curtain:       template.HTML(curtain),
		NavActive:     "repos",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs:        buildRepoCrumbs(owner, repo, ref, path, true),
	})
}

func (s *Server) renderFileView(w http.ResponseWriter, owner, repo, ref, path, fullName string, file *ghapi.RepoFile, repoInfo *ghapi.RepoInfo, sess *auth.Session, theme string) {
	var buf bytes.Buffer

	buf.WriteString(`<div class="phui-object-box">`)
	buf.WriteString(`<div class="phui-header-shell" style="display:flex;align-items:center;justify-content:space-between">`)
	buf.WriteString(`<div class="phui-header-view">`)
	icon := diff.FileIcon(file.Name)
	fmt.Fprintf(&buf, `<h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa %s"></span>%s`, icon, template.HTMLEscapeString(file.Name))
	fmt.Fprintf(&buf, ` <span style="font-weight:400;font-size:12px;color:#6b748c">%s</span></h1>`, formatSize(file.Size))
	buf.WriteString(`</div>`)
	if file.HTMLURL != "" {
		fmt.Fprintf(&buf, `<a href="%s" target="_blank" style="display:inline-flex;align-items:center;gap:4px;padding:4px 10px;font-size:12px;color:#6b748c;border:1px solid #c7ccd9;border-radius:3px;text-decoration:none;white-space:nowrap"><span class="phui-icon-view phui-font-fa fa-github" style="font-size:14px"></span>GitHub</a>`,
			template.HTMLEscapeString(file.HTMLURL))
	}
	buf.WriteString(`</div>`)

	ext := strings.ToLower(filepath.Ext(file.Name))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico", ".bmp":
		// Render image preview via GitHub raw URL
		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, path)
		buf.WriteString(`<div style="padding:24px;text-align:center;background:#f7f7f9">`)
		fmt.Fprintf(&buf, `<img src="%s" style="max-width:100%%;border:1px solid #c7ccd9;border-radius:3px" alt="%s">`,
			template.HTMLEscapeString(rawURL), template.HTMLEscapeString(file.Name))
		buf.WriteString(`</div>`)
	default:
		// Syntax-highlighted source
		lines := strings.Split(file.Content, "\n")
		highlighted := diff.HighlightLines(file.Name, lines)

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
	}
	buf.WriteString(`</div>`)

	// Curtain with repo info
	curtain := buildRepoCurtain(repoInfo)

	templates.RenderPage(w, templates.PageData{
		Title:         file.Name + " - Diffusion",
		Theme:         theme,
		Content:       template.HTML(buf.String()),
		Curtain:       template.HTML(curtain),
		NavActive:     "repos",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs:        buildRepoCrumbs(owner, repo, ref, path, false),
	})
}


func buildRepoCrumbs(owner, repo, ref, path string, isDir bool) []templates.Crumb {
	crumbs := []templates.Crumb{
		{Name: "Home", Href: "/"},
		{Name: "Repositories"},
	}

	repoHref := fmt.Sprintf("/repo/%s/%s?ref=%s", owner, repo, template.URLQueryEscaper(ref))
	if path == "" {
		crumbs = append(crumbs, templates.Crumb{Name: owner + "/" + repo})
		return crumbs
	}

	crumbs = append(crumbs, templates.Crumb{Name: owner + "/" + repo, Href: repoHref})

	parts := strings.Split(path, "/")
	for i, part := range parts {
		isLast := i == len(parts)-1
		if isLast {
			crumbs = append(crumbs, templates.Crumb{Name: part})
		} else {
			subPath := strings.Join(parts[:i+1], "/")
			href := fmt.Sprintf("/repo/%s/%s/%s?ref=%s", owner, repo, subPath, template.URLQueryEscaper(ref))
			crumbs = append(crumbs, templates.Crumb{Name: part, Href: href})
		}
	}

	return crumbs
}

func buildRepoCurtain(info *ghapi.RepoInfo) string {
	var c strings.Builder
	c.WriteString(`<div class="phui-box phui-box-border phui-object-box phui-curtain-view">`)

	// Actions
	c.WriteString(`<div class="phui-curtain-panel">`)
	c.WriteString(`<div class="phui-curtain-panel-header">Actions</div>`)
	c.WriteString(`<div class="phui-curtain-panel-body">`)
	c.WriteString(`<ul class="phabricator-action-list-view">`)
	fmt.Fprintf(&c, `<li class="phabricator-action-view action-has-icon"><a href="%s" target="_blank" class="phabricator-action-view-item"><span class="phabricator-action-view-icon phui-icon-view phui-font-fa fa-github"></span>View on GitHub</a></li>`,
		template.HTMLEscapeString(info.HTMLURL))
	c.WriteString(`</ul>`)
	c.WriteString(`</div></div>`)

	// Details as property list
	c.WriteString(`<div class="phui-curtain-panel">`)
	c.WriteString(`<div class="phui-curtain-panel-header">Details</div>`)
	c.WriteString(`<div class="phui-curtain-panel-body">`)
	c.WriteString(`<table class="phui-property-list-view" style="width:100%">`)

	visLabel := "Public"
	visIcon := "fa-globe"
	visShade := "blue"
	if info.Private {
		visLabel = "Private"
		visIcon = "fa-lock"
		visShade = "orange"
	}
	fmt.Fprintf(&c, `<tr><th style="color:#6b748c;font-weight:normal;padding:2px 8px 2px 0;white-space:nowrap;font-size:13px">Visibility</th><td style="padding:2px 0;font-size:13px"><span class="phui-tag-view phui-tag-shade-%s phui-tag-type-shade"><span class="phui-tag-core"><span class="phui-icon-view phui-font-fa %s mrs"></span>%s</span></span></td></tr>`, visShade, visIcon, visLabel)
	fmt.Fprintf(&c, `<tr><th style="color:#6b748c;font-weight:normal;padding:2px 8px 2px 0;white-space:nowrap;font-size:13px">Stars</th><td style="padding:2px 0;font-size:13px">%d</td></tr>`, info.Stars)
	fmt.Fprintf(&c, `<tr><th style="color:#6b748c;font-weight:normal;padding:2px 8px 2px 0;white-space:nowrap;font-size:13px">Forks</th><td style="padding:2px 0;font-size:13px">%d</td></tr>`, info.Forks)

	if info.Description != "" {
		fmt.Fprintf(&c, `<tr><th style="color:#6b748c;font-weight:normal;padding:2px 8px 2px 0;white-space:nowrap;font-size:13px;vertical-align:top">About</th><td style="padding:2px 0;font-size:13px">%s</td></tr>`, template.HTMLEscapeString(info.Description))
	}

	c.WriteString(`</table>`)
	c.WriteString(`</div></div>`)

	c.WriteString(`</div>`)
	return c.String()
}

func formatSize(bytes int) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%d B", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}
