package store

import "time"

func (s *Store) SaveSession(token []byte, userID string, expiresAt time.Time, userAgent string) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at, user_agent) VALUES (?, ?, ?, ?)",
		token, userID, expiresAt.Unix(), userAgent)
	return err
}

func (s *Store) DeleteSession(token []byte) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

func (s *Store) PurgeExpiredSessions() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now().Unix())
	return err
}
