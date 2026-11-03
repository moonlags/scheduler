package web

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

const (
	sessionCookie = "sid"
	sessionTTL    = 30 * 24 * time.Hour
)

// startSession creates a session for userID and sets the session cookie. Only
// the SHA-256 hash of the token is stored in the database.
func (s *Server) startSession(w http.ResponseWriter, r *http.Request, userID string) error {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	if err := s.store.SaveSession(hash(token), userID, time.Now().Add(sessionTTL), r.UserAgent()); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isHTTPS(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	return nil
}

func (s *Server) deleteSession(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		if err := s.store.DeleteSession(hash(c.Value)); err != nil {
			log.Println("delete session:", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isHTTPS(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func hash(data string) []byte {
	sum := sha256.Sum256([]byte(data))
	return sum[:]
}

// isHTTPS reports whether the request reached the server over HTTPS. Cookies
// marked Secure are not sent over plain HTTP, which would break local
// development.
func isHTTPS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}
