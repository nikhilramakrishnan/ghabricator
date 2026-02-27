package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	gh "github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
	oauthgithub "golang.org/x/oauth2/github"
)

type contextKey string

const (
	ctxSession    contextKey = "session"
	ctxGHClient   contextKey = "gh_client"
	stateCookieNm string     = "oauth_state"
)

type OAuthHandler struct {
	config *oauth2.Config
	store  *SessionStore
}

// Store returns the underlying session store.
func (h *OAuthHandler) Store() *SessionStore { return h.store }

func NewOAuthHandler(store *SessionStore) *OAuthHandler {
	return &OAuthHandler{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			Scopes:       []string{"repo", "gist"},
			Endpoint:     oauthgithub.Endpoint,
		},
		store: store,
	}
}

// HandleLogin redirects to GitHub's OAuth authorize page.
func (h *OAuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := randomState()
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieNm,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})
	// Remember the origin so we redirect back after OAuth callback.
	if origin := r.Header.Get("Origin"); origin != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_origin",
			Value:    origin,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})
	} else if ref := r.Referer(); ref != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_origin",
			Value:    ref,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})
	}
	url := h.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleCallback exchanges the OAuth code for a token and creates a session.
func (h *OAuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state
	stateCookie, err := r.Cookie(stateCookieNm)
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   stateCookieNm,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code parameter", http.StatusBadRequest)
		return
	}

	token, err := h.config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("oauth exchange error: %v", err)
		http.Error(w, "oauth exchange failed", http.StatusInternalServerError)
		return
	}

	// Fetch GitHub user info
	client := gh.NewClient(h.config.Client(context.Background(), token))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		log.Printf("github user fetch error: %v", err)
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	sess := h.store.Create(token, user.GetLogin(), user.GetAvatarURL())
	h.store.SetCookie(w, sess.ID)

	// Redirect back to the frontend origin if we saved one during login.
	redirect := "/"
	if c, err := r.Cookie("oauth_origin"); err == nil && c.Value != "" {
		redirect = c.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "oauth_origin",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// HandleLogout clears the session.
func (h *OAuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sess := h.store.GetFromRequest(r)
	if sess != nil {
		h.store.Delete(sess.ID)
	}
	h.store.ClearCookie(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// RequireAuth is middleware that ensures the request has a valid session.
// It stores the session and a GitHub client in the request context.
func (h *OAuthHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := h.store.GetFromRequest(r)
		if sess == nil {
			http.Redirect(w, r, "/auth/github", http.StatusTemporaryRedirect)
			return
		}
		client := gh.NewClient(h.config.Client(r.Context(), sess.Token))
		ctx := context.WithValue(r.Context(), ctxSession, sess)
		ctx = context.WithValue(ctx, ctxGHClient, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SessionFromContext retrieves the session from the request context.
func SessionFromContext(ctx context.Context) *Session {
	sess, _ := ctx.Value(ctxSession).(*Session)
	return sess
}

// GitHubClientFromContext retrieves the GitHub client from the request context.
func GitHubClientFromContext(ctx context.Context) *gh.Client {
	client, _ := ctx.Value(ctxGHClient).(*gh.Client)
	return client
}

func randomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
