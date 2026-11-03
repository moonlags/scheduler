package store

import (
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()

	st, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

func TestEnsureAdminAndSessions(t *testing.T) {
	st := newTestStore(t)

	if err := st.EnsureAdmin("admin@example.com", "secret"); err != nil {
		t.Fatalf("ensure admin: %v", err)
	}

	u, err := st.UserByEmail("admin@example.com")
	if err != nil {
		t.Fatalf("user by email: %v", err)
	}
	if !u.IsAdmin {
		t.Fatal("bootstrap user is not an admin")
	}

	token := []byte("test-token")
	if err := st.SaveSession(token, u.ID, time.Now().Add(time.Hour), "test-agent"); err != nil {
		t.Fatalf("save session: %v", err)
	}

	got, err := st.UserBySession(token)
	if err != nil {
		t.Fatalf("user by session: %v", err)
	}
	if got.ID != u.ID {
		t.Fatalf("user by session returned %q, want %q", got.ID, u.ID)
	}

	if err := st.DeleteSession(token); err != nil {
		t.Fatalf("delete session: %v", err)
	}
	if _, err := st.UserBySession(token); err == nil {
		t.Fatal("deleted session still resolves a user")
	}
}

func TestEnsureAdminUpdatesExistingAccount(t *testing.T) {
	st := newTestStore(t)

	if err := st.EnsureAdmin("admin@example.com", "first"); err != nil {
		t.Fatalf("ensure admin: %v", err)
	}
	if err := st.EnsureAdmin("admin@example.com", "second"); err != nil {
		t.Fatalf("ensure admin again: %v", err)
	}

	var count int
	if err := st.db.Get(&count, "SELECT COUNT(*) FROM users"); err != nil {
		t.Fatalf("count users: %v", err)
	}
	if count != 1 {
		t.Fatalf("got %d users, want 1", count)
	}

	u, err := st.UserByEmail("admin@example.com")
	if err != nil {
		t.Fatalf("user by email: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("second")); err != nil {
		t.Fatalf("password was not updated: %v", err)
	}
}
