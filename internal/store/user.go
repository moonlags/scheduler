package store

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	Email          string `db:"email"`
	Password       string `db:"password"`
	ProfilePicture string `db:"profile_picture"`
	IsAdmin        bool   `db:"is_admin"`
}

const userColumns = "id, name, email, password, profile_picture, is_admin"

func (s *Store) InsertUser(name, email, password, profilePicture string, isAdmin bool) (User, error) {
	u := User{
		ID:             ulid.Make().String(),
		Name:           name,
		Email:          email,
		Password:       password,
		ProfilePicture: profilePicture,
		IsAdmin:        isAdmin,
	}

	if _, err := s.db.Exec(
		"INSERT INTO users (id, name, email, password, profile_picture, is_admin) VALUES (?, ?, ?, ?, ?, ?)",
		u.ID, u.Name, u.Email, u.Password, u.ProfilePicture, u.IsAdmin); err != nil {
		return User{}, fmt.Errorf("insert user: %w", err)
	}

	return u, nil
}

func (s *Store) UserByEmail(email string) (User, error) {
	var u User
	err := s.db.Get(&u, "SELECT "+userColumns+" FROM users WHERE email = ?", email)
	return u, err
}

// UserBySession returns the user owning a session token that has not expired.
func (s *Store) UserBySession(token []byte) (User, error) {
	const q = `
		SELECT u.id, u.name, u.email, u.password, u.profile_picture, u.is_admin
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token = ? AND s.expires_at > ?`

	var u User
	err := s.db.Get(&u, q, token, time.Now().Unix())
	return u, err
}

// EnsureAdmin creates the bootstrap administrator or updates its password if
// the account already exists. The admin account is identified by email.
func (s *Store) EnsureAdmin(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	const q = `
		INSERT INTO users (id, name, email, password, profile_picture, is_admin)
		VALUES (?, 'admin', ?, ?, 'default', TRUE)
		ON CONFLICT (email) DO UPDATE SET
			password = excluded.password,
			is_admin = TRUE`

	if _, err := s.db.Exec(q, ulid.Make().String(), email, hash); err != nil {
		return fmt.Errorf("upsert admin: %w", err)
	}
	return nil
}
