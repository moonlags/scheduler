package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"
)

func (s *Server) saveSession(hashedToken []byte, userID string, expiresAt time.Time, userAgent string) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at, userAgent) VALUES (?, ?, ?, ?)",
		hashedToken, userID, expiresAt.Unix())
	return err
}

func (s *Server) deleteSession(hashedToken string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = ?", hashedToken)
	return err
}

func (s *Server) startSession(w http.ResponseWriter, r *http.Request, userID string) error {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	if err := s.saveSession(hash(token), userID, time.Now().Add(30*24*time.Hour), r.UserAgent()); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 3600,
	})
	return nil
}

func (s *Server) purgeExpiredSessions() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now().Unix())
	return err
}

func hash(data string) []byte {
	sum := sha256.Sum256([]byte(data))
	return sum[:]
}
