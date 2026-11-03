package web

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/moonlags/ikeaScheduler/internal/store"
)

type ctxKey struct{}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "login.html", nil)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" {
		http.Error(w, "email is not optional", http.StatusUnprocessableEntity)
		return
	}
	if password == "" {
		http.Error(w, "password is not optional", http.StatusUnprocessableEntity)
		return
	}

	u, err := s.store.UserByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "email or password is not correct", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("handleLogin:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		http.Error(w, "email or password is not correct", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("handleLogin:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := s.startSession(w, r, u.ID); err != nil {
		log.Println("handleLogin:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	redirect(w, r, "/dashboard")
}

func (s *Server) handleLogoutPage(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "logout.html", nil)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.deleteSession(w, r)
	redirect(w, r, "/login")
}

// withUser loads the logged-in user, if any, into the request context.
func (s *Server) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(sessionCookie); err == nil {
			if u, err := s.store.UserBySession(hash(c.Value)); err == nil {
				r = r.WithContext(context.WithValue(r.Context(), ctxKey{}, &u))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userFrom(r.Context()) == nil {
			redirect(w, r, "/login")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := userFrom(r.Context())
		if u == nil {
			redirect(w, r, "/login")
			return
		}
		if !u.IsAdmin {
			redirect(w, r, "/dashboard")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func userFrom(ctx context.Context) *store.User {
	u, _ := ctx.Value(ctxKey{}).(*store.User)
	return u
}
