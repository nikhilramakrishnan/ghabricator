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

	// Branch selector
	writeBranchSelector(&buf, owner, repo, ref, path)

	buf.WriteString(`<div class="phui-box phui-box-border phui-object-box">`)
	buf.WriteString(`<div class="phui-header-shell">`)
	buf.WriteString(`<div class="phui-header-view">`)
	buf.WriteString(`<span class="phui-header-icon phui-icon-view phui-font-fa fa-folder-open"></span>`)
	fmt.Fprintf(&buf, `<h1 class="phui-header-header">%s</h1>`, fullName)
	if path != "" {
		fmt.Fprintf(&buf, `<span style="margin-left:8px;color:#6b748c">/ %s</span>`, template.HTMLEscapeString(path))
	}
	buf.WriteString(`</div>`)
	buf.WriteString(`</div>`)

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
		HeaderTitle:   template.HTML(fullName),
		HeaderIcon:    "fa-folder-open",
		Content:       template.HTML(buf.String()),
		Curtain:       template.HTML(curtain),
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs:        buildRepoCrumbs(owner, repo, ref, path, true),
	})
}

func (s *Server) renderFileView(w http.ResponseWriter, owner, repo, ref, path, fullName string, file *ghapi.RepoFile, repoInfo *ghapi.RepoInfo, sess *auth.Session, theme string) {
	var buf bytes.Buffer

	// Branch selector
	writeBranchSelector(&buf, owner, repo, ref, path)

	// File header
	buf.WriteString(`<div class="phui-object-box">`)
	buf.WriteString(`<div class="phui-header-shell" style="display:flex;align-items:center;justify-content:space-between">`)
	buf.WriteString(`<div class="phui-header-view">`)
	icon := diff.FileIcon(file.Name)
	fmt.Fprintf(&buf, `<span class="phui-header-icon phui-icon-view phui-font-fa %s"></span>`, icon)
	fmt.Fprintf(&buf, `<h1 class="phui-header-header">%s</h1>`, template.HTMLEscapeString(file.Name))
	fmt.Fprintf(&buf, `<span style="margin-left:12px;color:#6b748c;font-size:13px">%s</span>`, formatSize(file.Size))
	buf.WriteString(`</div>`)

	// Raw link
	if file.HTMLURL != "" {
		fmt.Fprintf(&buf, `<a href="%s" target="_blank" class="phui-button-view button-grey" style="margin-right:8px"><span class="phui-icon-view phui-font-fa fa-github mrs"></span>View on GitHub</a>`,
			template.HTMLEscapeString(file.HTMLURL))
	}
	buf.WriteString(`</div>`)

	// Syntax-highlighted source
	lines := strings.Split(file.Content, "\n")
	highlighted := diff.HighlightLines(file.Name, lines)

	buf.WriteString(`<div class="phabricator-source-code-container">`)
	buf.WriteString(`<table class="phabricator-source-code-view remarkup-code PhabricatorMonospaced">`)
	for i, hl := range highlighted {
		lineNum := i + 1
		buf.WriteString(`<tr>`)
		fmt.Fprintf(&buf, `<th class="phabricator-source-line"><span>%d</span></th>`, lineNum)
		fmt.Fprintf(&buf, `<td class="phabricator-source-code">%s</td>`, hl)
		buf.WriteString(`</tr>`)
	}
	buf.WriteString(`</table>`)
	buf.WriteString(`</div>`)
	buf.WriteString(`</div>`)

	// Curtain with repo info
	curtain := buildRepoCurtain(repoInfo)

	templates.RenderPage(w, templates.PageData{
		Title:         file.Name + " - Diffusion",
		Theme:         theme,
		HeaderTitle:   template.HTML(fullName),
		HeaderIcon:    icon,
		Content:       template.HTML(buf.String()),
		Curtain:       template.HTML(curtain),
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs:        buildRepoCrumbs(owner, repo, ref, path, false),
	})
}

func writeBranchSelector(buf *bytes.Buffer, owner, repo, ref, path string) {
	buf.WriteString(`<div style="margin-bottom:12px;display:flex;align-items:center;gap:8px">`)
	buf.WriteString(`<span class="phui-icon-view phui-font-fa fa-code-fork" style="color:#6b748c"></span>`)
	buf.WriteString(`<form method="GET" style="display:inline">`)
	// Preserve path — form action is the current page
	action := fmt.Sprintf("/repo/%s/%s", owner, repo)
	if path != "" {
		action += "/" + path
	}
	fmt.Fprintf(buf, `<input type="text" name="ref" value="%s" placeholder="branch or tag" style="padding:4px 8px;border:1px solid #c7ccd9;border-radius:3px;font-size:13px;width:200px">`, template.HTMLEscapeString(ref))
	buf.WriteString(` <button type="submit" class="phui-button-view button-grey" style="font-size:12px;padding:4px 10px">Switch</button>`)
	buf.WriteString(`</form>`)
	buf.WriteString(`</div>`)
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
	var curtain strings.Builder
	curtain.WriteString(`<div class="phui-box phui-box-border phui-object-box phui-curtain-view">`)

	// Repo details
	curtain.WriteString(`<div class="phui-curtain-panel">`)
	curtain.WriteString(`<div class="phui-curtain-panel-header">Repository</div>`)
	curtain.WriteString(`<div class="phui-curtain-panel-body">`)
	fmt.Fprintf(&curtain, `<div><span class="phui-icon-view phui-font-fa fa-github mrs"></span>%s</div>`, template.HTMLEscapeString(info.FullName))
	if info.Description != "" {
		fmt.Fprintf(&curtain, `<div style="margin-top:4px;color:#6b748c;font-size:13px">%s</div>`, template.HTMLEscapeString(info.Description))
	}
	visIcon := "fa-globe"
	visLabel := "Public"
	if info.Private {
		visIcon = "fa-lock"
		visLabel = "Private"
	}
	fmt.Fprintf(&curtain, `<div style="margin-top:8px"><span class="phui-icon-view phui-font-fa %s mrs"></span>%s</div>`, visIcon, visLabel)
	fmt.Fprintf(&curtain, `<div><span class="phui-icon-view phui-font-fa fa-star mrs"></span>%d stars</div>`, info.Stars)
	fmt.Fprintf(&curtain, `<div><span class="phui-icon-view phui-font-fa fa-code-fork mrs"></span>%d forks</div>`, info.Forks)
	curtain.WriteString(`</div></div>`)

	// Actions
	curtain.WriteString(`<div class="phui-curtain-panel">`)
	curtain.WriteString(`<div class="phui-curtain-panel-header">Actions</div>`)
	curtain.WriteString(`<div class="phui-curtain-panel-body">`)
	curtain.WriteString(`<ul class="phabricator-action-list-view">`)
	fmt.Fprintf(&curtain, `<li class="phabricator-action-view action-has-icon"><a href="%s" target="_blank" class="phabricator-action-view-item"><span class="phabricator-action-view-icon phui-icon-view phui-font-fa fa-github"></span>View on GitHub</a></li>`,
		template.HTMLEscapeString(info.HTMLURL))
	curtain.WriteString(`</ul>`)
	curtain.WriteString(`</div></div>`)

	curtain.WriteString(`</div>`)
	return curtain.String()
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
