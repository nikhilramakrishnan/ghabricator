package server

import (
	"net/http"
	"os"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/herald"
)

type Server struct {
	mux    *http.ServeMux
	auth   *auth.OAuthHandler
	herald *herald.Store
}

func New() (*Server, error) {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}
	store := auth.NewSessionStore(secret)

	s := &Server{
		mux:    http.NewServeMux(),
		auth:   auth.NewOAuthHandler(store),
		herald: herald.NewStore(),
	}
	s.routes()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	// Auth routes
	s.mux.HandleFunc("GET /auth/github", s.auth.HandleLogin)
	s.mux.HandleFunc("GET /auth/callback", s.auth.HandleCallback)
	s.mux.HandleFunc("GET /auth/logout", s.auth.HandleLogout)
	s.mux.HandleFunc("GET /api/auth/github", s.auth.HandleLogin)
	s.mux.HandleFunc("GET /api/auth/callback", s.auth.HandleCallback)
	s.mux.HandleFunc("GET /api/auth/logout", s.auth.HandleLogout)
	s.mux.HandleFunc("GET /api/auth/me", s.handleAPIAuthMe)

	// Dashboard
	s.mux.Handle("GET /api/dashboard", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIDashboard)))

	// PR
	s.mux.Handle("GET /api/pr/{owner}/{repo}/{number}", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIPR)))

	// PR compare (diff between two commits)
	s.mux.Handle("GET /api/pr/{owner}/{repo}/{number}/compare", s.auth.RequireAuth(http.HandlerFunc(s.handleAPICompare)))

	// Inline comments
	s.mux.Handle("POST /api/v2/inline", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIInline)))

	// Review / merge / close
	s.mux.Handle("POST /api/v2/review", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIReview)))
	s.mux.Handle("POST /api/v2/merge", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIMerge)))
	s.mux.Handle("POST /api/v2/close", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIClose)))

	// Edit PR / comments
	s.mux.Handle("POST /api/v2/edit-pr", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIEditPR)))
	s.mux.Handle("POST /api/v2/edit-comment", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIEditComment)))

	// Reactions
	s.mux.Handle("POST /api/v2/reaction", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIReaction)))

	// Repos
	s.mux.Handle("GET /api/repos", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIRepos)))
	s.mux.Handle("GET /api/repo/{owner}/{repo}/info", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIRepoInfo)))
	s.mux.Handle("GET /api/repo/{owner}/{repo}/tree", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIRepoTree)))
	s.mux.Handle("GET /api/repo/{owner}/{repo}/file", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIRepoFile)))

	// Paste
	s.mux.Handle("GET /api/paste", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIPasteList)))
	s.mux.Handle("GET /api/paste/{id}", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIPasteView)))
	s.mux.Handle("POST /api/paste", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIPasteCreate)))

	// Herald
	s.mux.Handle("GET /api/herald", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIHeraldList)))
	s.mux.Handle("GET /api/herald/{id}", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIHeraldGet)))
	s.mux.Handle("POST /api/herald", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIHeraldSave)))
	s.mux.Handle("DELETE /api/herald/{id}", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIHeraldDelete)))

	// Search
	s.mux.Handle("GET /api/search", s.auth.RequireAuth(http.HandlerFunc(s.handleAPISearch)))

	// Actions (workflow runs)
	s.mux.Handle("GET /api/actions/runs", s.auth.RequireAuth(http.HandlerFunc(s.handleAPIWorkflowRuns)))
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



