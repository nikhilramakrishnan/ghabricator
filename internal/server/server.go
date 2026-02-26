package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/nikhilr/ghabricator/internal/assets"
	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/herald"
	"github.com/nikhilr/ghabricator/internal/templates"
)

type Server struct {
	mux      *http.ServeMux
	auth     *auth.OAuthHandler
	assets   *assets.Server
	herald   *herald.Store
	repoRoot string
}

func New(repoRoot string) (*Server, error) {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}
	store := auth.NewSessionStore(secret)

	assetSrv, err := assets.NewServer(repoRoot)
	if err != nil {
		return nil, err
	}

	s := &Server{
		mux:      http.NewServeMux(),
		auth:     auth.NewOAuthHandler(store),
		assets:   assetSrv,
		herald:   herald.NewStore(),
		repoRoot: repoRoot,
	}
	s.routes()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	// Public routes
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("GET /res/", s.handleAssets)
	s.mux.HandleFunc("GET /rsrc/", s.handleRsrc)
	s.mux.HandleFunc("GET /api/theme", s.handleTheme)

	// Auth routes
	s.mux.HandleFunc("GET /auth/github", s.auth.HandleLogin)
	s.mux.HandleFunc("GET /auth/callback", s.auth.HandleCallback)
	s.mux.HandleFunc("GET /auth/logout", s.auth.HandleLogout)

	// Protected routes (require auth)
	s.mux.Handle("GET /dashboard", s.auth.RequireAuth(http.HandlerFunc(s.handleDashboard)))
	s.mux.Handle("GET /pr/{owner}/{repo}/{number}", s.auth.RequireAuth(http.HandlerFunc(s.handlePR)))
	s.mux.Handle("POST /api/inline", s.auth.RequireAuth(http.HandlerFunc(s.handleInline)))
	s.mux.Handle("POST /api/review", s.auth.RequireAuth(http.HandlerFunc(s.handleReview)))

	// Paste routes
	s.mux.Handle("GET /paste", s.auth.RequireAuth(http.HandlerFunc(s.handlePasteList)))
	s.mux.Handle("GET /paste/new", s.auth.RequireAuth(http.HandlerFunc(s.handlePasteNew)))
	s.mux.Handle("POST /paste/save", s.auth.RequireAuth(http.HandlerFunc(s.handlePasteSave)))
	s.mux.Handle("GET /paste/{id}", s.auth.RequireAuth(http.HandlerFunc(s.handlePasteView)))

	// Herald routes
	s.mux.Handle("GET /herald", s.auth.RequireAuth(http.HandlerFunc(s.handleHeraldList)))
	s.mux.Handle("GET /herald/new", s.auth.RequireAuth(http.HandlerFunc(s.handleHeraldNew)))
	s.mux.Handle("POST /herald/save", s.auth.RequireAuth(http.HandlerFunc(s.handleHeraldSave)))
	s.mux.Handle("GET /herald/{id}", s.auth.RequireAuth(http.HandlerFunc(s.handleHeraldView)))
	s.mux.Handle("GET /herald/{id}/delete", s.auth.RequireAuth(http.HandlerFunc(s.handleHeraldDelete)))

	// Repository browser
	s.mux.Handle("GET /repos", s.auth.RequireAuth(http.HandlerFunc(s.handleRepoList)))
	s.mux.Handle("GET /repo/{owner}/{repo}/{path...}", s.auth.RequireAuth(http.HandlerFunc(s.handleRepoView)))
	s.mux.Handle("GET /repo/{owner}/{repo}", s.auth.RequireAuth(http.HandlerFunc(s.handleRepoView)))

	// Search
	s.mux.Handle("GET /search", s.auth.RequireAuth(http.HandlerFunc(s.handleSearch)))

	// PR actions (merge / close / reopen)
	s.mux.Handle("POST /api/merge", s.auth.RequireAuth(http.HandlerFunc(s.handleMerge)))
	s.mux.Handle("POST /api/close", s.auth.RequireAuth(http.HandlerFunc(s.handleClose)))

	// Diff context expansion
	s.mux.Handle("GET /api/context", s.auth.RequireAuth(http.HandlerFunc(s.handleContext)))

	// JSON API routes (SvelteKit frontend)
	s.mux.HandleFunc("GET /api/auth/me", s.handleAPIAuthMe)
	s.mux.Handle("GET /api/dashboard", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIDashboard)))
	s.mux.Handle("GET /api/pr/{owner}/{repo}/{number}", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIPR)))
}

func (s *Server) handleAPIAuthMe(w http.ResponseWriter, r *http.Request) {
	sess := s.auth.Store().GetFromRequest(r)
	if sess == nil {
		jsonError(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	jsonOK(w, map[string]string{
		"login":     sess.Login,
		"avatarURL": sess.AvatarURL,
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// If authenticated, redirect to dashboard
	sess := s.auth.Store().GetFromRequest(r)
	if sess != nil {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}

	content := `<div class="phui-box phui-box-border phui-object-box" style="text-align:center;padding:40px 20px">
  <div style="margin-bottom:24px">
    <span class="phui-icon-view phui-font-fa fa-code-fork" style="font-size:48px;color:#6b748c"></span>
  </div>
  <h2 style="margin-bottom:12px;font-size:24px">Welcome to Ghabricator</h2>
  <p style="color:#6b748c;font-size:15px;margin-bottom:24px;max-width:480px;margin-left:auto;margin-right:auto">
    A modern code review experience powered by GitHub. Review pull requests with Phabricator's legendary diff viewer, inline comments, and Herald automation.
  </p>
  <a href="/auth/github" class="phui-button-view button-green" style="font-size:16px;padding:10px 28px;display:inline-flex;align-items:center;gap:8px">
    <span class="phui-icon-view phui-font-fa fa-github"></span> Sign in with GitHub
  </a>
  <div style="margin-top:32px;color:#92969d;font-size:13px">
    <div style="display:flex;justify-content:center;gap:32px;flex-wrap:wrap">
      <span><span class="phui-icon-view phui-font-fa fa-search mrs"></span> Side-by-side diffs</span>
      <span><span class="phui-icon-view phui-font-fa fa-commenting-o mrs"></span> Inline comments</span>
      <span><span class="phui-icon-view phui-font-fa fa-bullhorn mrs"></span> Herald rules</span>
      <span><span class="phui-icon-view phui-font-fa fa-clipboard mrs"></span> Paste bin</span>
    </div>
  </div>
</div>`

	theme := templates.ThemeFromRequest(r)
	templates.RenderPage(w, templates.PageData{
		Title:   "Ghabricator",
		Theme:   theme,
		Content: template.HTML(content),
	})
}

func (s *Server) handleTheme(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	maxAge := 365 * 24 * 60 * 60
	value := ""
	if mode == "dark" {
		value = "dark"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	})
	ref := r.Referer()
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusTemporaryRedirect)
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	s.assets.ServeHTTP(w, r)
}

func (s *Server) renderError(w http.ResponseWriter, r *http.Request, title string, message string, code int) {
	w.WriteHeader(code)
	theme := templates.ThemeFromRequest(r)

	var icon, shade string
	switch {
	case code >= 500:
		icon = "fa-exclamation-triangle"
		shade = "red"
	case code == 404:
		icon = "fa-question-circle"
		shade = "blue"
	case code == 403:
		icon = "fa-lock"
		shade = "yellow"
	default:
		icon = "fa-info-circle"
		shade = "blue"
	}

	content := fmt.Sprintf(
		`<div class="phui-box phui-box-border phui-object-box">
        <div class="phui-info-view phui-info-severity-%s">
            <span class="phui-icon-view phui-font-fa %s mrs"></span>
            <strong>%s</strong>
            <p style="margin-top:8px">%s</p>
            <p style="margin-top:12px"><a href="/" class="phui-button-view button-grey" style="text-decoration:none">
                <span class="phui-icon-view phui-font-fa fa-home mrs"></span>Return Home</a></p>
        </div></div>`,
		shade, icon,
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(message),
	)

	templates.RenderPage(w, templates.PageData{
		Title:   title,
		Theme:   theme,
		Content: template.HTML(content),
		Crumbs:  []templates.Crumb{{Name: "Home", Href: "/"}, {Name: title}},
	})
}

func (s *Server) handleRsrc(w http.ResponseWriter, r *http.Request) {
	// Phabricator's font-face CSS references /rsrc/ paths directly.
	// Rewrite to serve from the webroot.
	r.URL.Path = "/res/" + strings.TrimPrefix(r.URL.Path, "/rsrc/")
	s.assets.ServeHTTP(w, r)
}


