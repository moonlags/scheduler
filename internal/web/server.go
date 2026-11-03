// Package web contains the HTTP layer: routing, middleware, handlers and
// server-rendered templates.
package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/moonlags/ikeaScheduler/internal/store"
)

type Server struct {
	store *store.Store
	mux   *chi.Mux
	pages map[string]*template.Template
	frags *template.Template
}

func New(st *store.Store) (*Server, error) {
	s := &Server{
		store: st,
		mux:   chi.NewRouter(),
	}

	if err := s.loadTemplates(); err != nil {
		return nil, err
	}
	if err := s.routes(); err != nil {
		return nil, err
	}
	return s, nil
}

// Handler exposes the router for tests and embedding.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Run(addr string) error {
	go s.purgeSessions()

	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() error {
	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		return fmt.Errorf("static files: %w", err)
	}

	cop := http.NewCrossOriginProtection()

	s.mux.Use(middleware.RequestID)
	s.mux.Use(httprate.LimitBy(64, 2*time.Minute, func(r *http.Request) (string, error) {
		return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
	}))
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
	s.mux.Use(cop.Handler)
	s.mux.Use(s.withUser)

	s.mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	s.mux.Get("/", s.handleIndexPage)
	s.mux.Get("/login", s.handleLoginPage)
	s.mux.Post("/login", s.handleLogin)

	s.mux.Group(func(r chi.Router) {
		r.Use(s.requireUser)

		r.Get("/logout", s.handleLogoutPage)
		r.Post("/logout", s.handleLogout)

		r.Get("/dashboard", s.handleDashboardPage)
	})

	// TODO: add routes for invite-based registration, profile management,
	// organizations and employees together with their handlers.

	return nil
}

func (s *Server) purgeSessions() {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for range t.C {
		if err := s.store.PurgeExpiredSessions(); err != nil {
			log.Println("purge sessions:", err)
		}
	}
}
