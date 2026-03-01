package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

// AuthHandler handles authentication via either OAuth or a static PAT.
type AuthHandler struct {
	// OAuth mode fields (nil in token mode).
	config *oauth2.Config
	store  *SessionStore

	// Token mode fields (nil in OAuth mode).
	tokenSession *Session
	tokenClient  *gh.Client
}

// Store returns the underlying session store (nil in token mode).
func (h *AuthHandler) Store() *SessionStore { return h.store }

// TokenSession returns the static session in token mode (nil in OAuth mode).
func (h *AuthHandler) TokenSession() *Session { return h.tokenSession }

// IsTokenMode returns true if using a static PAT instead of OAuth.
func (h *AuthHandler) IsTokenMode() bool { return h.tokenSession != nil }

// NewAuthHandler creates an auth handler. It auto-detects the mode:
//   - If GITHUB_TOKEN is set → token mode (no OAuth app needed)
//   - Otherwise → OAuth mode (requires GITHUB_CLIENT_ID + GITHUB_CLIENT_SECRET)
func NewAuthHandler(store *SessionStore) (*AuthHandler, error) {
	if pat := os.Getenv("GITHUB_TOKEN"); pat != "" {
		return newTokenHandler(pat)
	}
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("set GITHUB_TOKEN for PAT mode, or GITHUB_CLIENT_ID + GITHUB_CLIENT_SECRET for OAuth mode")
	}
	return &AuthHandler{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"repo", "gist"},
			Endpoint:     oauthgithub.Endpoint,
		},
		store: store,
	}, nil
}

func newTokenHandler(pat string) (*AuthHandler, error) {
	// Build a static client from the PAT.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: pat})
	client := gh.NewClient(oauth2.NewClient(ctx, ts))

	// Fetch the authenticated user to populate session info.
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("GITHUB_TOKEN invalid or cannot reach GitHub: %w", err)
	}

	sess := &Session{
		ID:        "token-mode",
		Token:     &oauth2.Token{AccessToken: pat},
		Login:     user.GetLogin(),
		AvatarURL: user.GetAvatarURL(),
	}

	log.Printf("Token mode: authenticated as %s", sess.Login)
	return &AuthHandler{
		tokenSession: sess,
		tokenClient:  client,
	}, nil
}

// HandleLogin redirects to GitHub's OAuth authorize page.
// In token mode, just redirects to /.
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if h.IsTokenMode() {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
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
// No-op in token mode.
func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if h.IsTokenMode() {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
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
// In token mode, just redirects to / (can't log out of a PAT).
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if h.IsTokenMode() {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sess := h.store.GetFromRequest(r)
	if sess != nil {
		h.store.Delete(sess.ID)
	}
	h.store.ClearCookie(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// RequireAuth is middleware that ensures the request has a valid session.
// It stores the session and a GitHub client in the request context.
func (h *AuthHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sess *Session
		var client *gh.Client

		if h.IsTokenMode() {
			sess = h.tokenSession
			client = h.tokenClient
		} else {
			sess = h.store.GetFromRequest(r)
			if sess == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"not authenticated"}`))
				return
			}
			client = gh.NewClient(h.config.Client(r.Context(), sess.Token))
		}

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
