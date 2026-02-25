package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	sessionCookieName = "phab_session"
	sessionTTL        = 24 * time.Hour
)

type Session struct {
	ID        string
	Token     *oauth2.Token
	Login     string
	AvatarURL string
	CreatedAt time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	secret   []byte
}

func NewSessionStore(secret string) *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
		secret:   []byte(secret),
	}
}

func (s *SessionStore) Create(token *oauth2.Token, login, avatarURL string) *Session {
	id := randomID(32)
	sess := &Session{
		ID:        id,
		Token:     token,
		Login:     login,
		AvatarURL: avatarURL,
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.sessions[id] = sess
	s.mu.Unlock()
	return sess
}

func (s *SessionStore) Get(id string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[id]
	if !ok {
		return nil
	}
	if time.Since(sess.CreatedAt) > sessionTTL {
		return nil
	}
	return sess
}

func (s *SessionStore) Delete(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}

func (s *SessionStore) SetCookie(w http.ResponseWriter, sessionID string) {
	sig := s.sign(sessionID)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID + "." + sig,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
}

func (s *SessionStore) GetFromRequest(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}
	id, sig := splitCookieValue(cookie.Value)
	if id == "" || !s.verify(id, sig) {
		return nil
	}
	return s.Get(id)
}

func (s *SessionStore) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (s *SessionStore) sign(data string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *SessionStore) verify(data, sig string) bool {
	expected := s.sign(data)
	return hmac.Equal([]byte(expected), []byte(sig))
}

func splitCookieValue(val string) (string, string) {
	for i := len(val) - 1; i >= 0; i-- {
		if val[i] == '.' {
			return val[:i], val[i+1:]
		}
	}
	return "", ""
}

func randomID(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
