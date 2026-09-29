package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

	u, err := s.getUserByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "email or password is not correct", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("handleLogin: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		http.Error(w, "email or password is not correct", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("handleLogin: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := s.startSession(w, r, u.ID); err != nil {
		log.Println("handleLogin: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/dashboard")
}

// func (s *Server) withUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if c, err := r.Cookie("sid"); err == nil {
// 			if u, err := s.store.UserBySession(r.Context(), hash(c.Value)); err == nil {
// 				r = r.WithContext(context.WithValue(r.Context(), ctxKey{}, u))
// 			}
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// /r.Group(func(r chi.Router) {
// 		r.Use(s.requireUser)
// 		r.Get("/profile", s.handleProfile)
// 		r.Post("/profile/avatar", s.handleAvatarUpload)
// 	})/ }
// r.Use(s.withUser)
// func (s *Server) requireUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if userFrom(r.Context()) == nil {
// 			if r.Header.Get("HX-Request") == "true" {
// 				w.Header().Set("HX-Redirect", "/login")
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
//
// func userFrom(ctx context.Context) *User {
// 	u, _ := ctx.Value(ctxKey{}).(*User)
// 	return u
// }
// go func() {
// 	t := time.NewTicker(time.Hour)
// 	defer t.Stop()
// 	for range t.C {
// 		if err := s.store.PurgeExpiredSessions(context.Background()); err != nil {
// 			log.Println("session purge:", err)
// 		}
// 	}
// }()
// // DELETE FROM sessions WHERE expires_at < now();
