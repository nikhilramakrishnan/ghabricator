package templates

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed *.html
var templateFS embed.FS

var pageTmpl = template.Must(template.ParseFS(templateFS, "page.html"))

// JavelinInit represents a Javelin initialization data tag.
type JavelinInit struct {
	Kind string
	Data template.HTML // Raw JSON string (must not be HTML-escaped)
}

// Crumb represents a breadcrumb item.
type Crumb struct {
	Name string
	Href string
}

// PageData holds all data needed to render the page shell.
type PageData struct {
	Title       string
	BodyClass   string
	Theme       string // "dark" or "" (default)
	Content     template.HTML
	Curtain     template.HTML
	FileTree    template.HTML // Left sidebar file tree (formation view)
	ExtraCSS    []string
	ExtraJS     []string
	JavelinData []JavelinInit
	Crumbs      []Crumb
	HeaderTitle template.HTML
	HeaderIcon  string
	InlineScript template.HTML // Inline JS injected before </body>
	// Nav state
	UserLogin     string
	UserAvatarURL string
}

// HasFileTree returns true if a file tree sidebar should be rendered.
func (d PageData) HasFileTree() bool {
	return d.FileTree != ""
}

// CSSPackagePath returns the core CSS path based on theme.
func (d PageData) CSSPackagePath() string {
	if d.Theme == "dark" {
		return "/res/pkg/dark/core.pkg.css"
	}
	return "/res/pkg/core.pkg.css"
}

// ThemeToggleIcon returns the FA icon for the theme toggle button.
func (d PageData) ThemeToggleIcon() string {
	if d.Theme == "dark" {
		return "fa-sun-o"
	}
	return "fa-moon-o"
}

// ThemeToggleMode returns the mode to switch to.
func (d PageData) ThemeToggleMode() string {
	if d.Theme == "dark" {
		return "default"
	}
	return "dark"
}

// ThemeFromRequest reads the theme cookie from a request.
func ThemeFromRequest(r *http.Request) string {
	c, err := r.Cookie("theme")
	if err != nil || c.Value != "dark" {
		return ""
	}
	return "dark"
}

// RenderPage writes the full page shell to the response.
func RenderPage(w http.ResponseWriter, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
