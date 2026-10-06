package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	cfg   Config
	mux   *chi.Mux
	db    *sqlx.DB
	pages map[string]*template.Template
	frags *template.Template
}

func newServer(cfg Config, db *sqlx.DB) (*Server, error) {
	server := &Server{
		mux: chi.NewRouter(),
		db:  db,
		cfg: cfg,
	}

	return server, nil
}

func (s *Server) run(addr string) error {
	log.Printf("running on address %s\n", addr)

	go func() {
		t := time.NewTicker(time.Hour)
		defer t.Stop()
		for range t.C {
			if err := s.purgeExpiredSessions(); err != nil {
				log.Println("session purge:", err)
			}
		}
	}()

	s.loadTemplates(os.DirFS("."))
	s.routes()
	s.initAdmin()

	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) initAdmin() {
	hash, err := bcrypt.GenerateFromPassword([]byte(s.cfg.adminPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("hash admin password: %v", err))
	}

	if _, err := s.db.Exec("INSERT INTO users (id, name, email, password, profile_picture, is_admin) VALUES (?, ?, ?, ?, 'default', TRUE) ON CONFLICT (id) DO UPDATE SET email = excluded.email, password = excluded.password",
		"1", "admin", s.cfg.adminEmail, hash); err != nil {
		panic(fmt.Sprintf("init admin: %v", err))
	}
}

func (s *Server) loadTemplates(fsys fs.FS) {
	pages, err := fs.Glob(fsys, "templates/pages/*.html")
	if err != nil {
		panic(fmt.Sprintf("read templates/pages/*.html: %v", err))
	}
	s.pages = make(map[string]*template.Template, len(pages))

	for _, page := range pages {
		tmpl, err := template.ParseFS(fsys, "templates/layout.html", "templates/fragments/*.html", page)
		if err != nil {
			panic(fmt.Sprintf("parse page %s: %v", page, err))
		}
		s.pages[path.Base(page)] = tmpl
	}

	s.frags = template.Must(template.ParseFS(fsys, "templates/fragments/*.html"))
}

func (s *Server) routes() {
	cop := http.NewCrossOriginProtection()

	s.mux.Use(middleware.RequestID)
	s.mux.Use(httprate.LimitBy(64, 2*time.Minute, func(r *http.Request) (string, error) {
		return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
	}))
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
	s.mux.Use(func(next http.Handler) http.Handler { return cop.Handler(next) })
	s.mux.Use(s.withUser)

	s.mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	s.mux.Get("/", s.handleIndexPage)

	s.mux.Get("/login", s.handleLoginPage)
	s.mux.Post("/login", s.handleLogin)

	s.mux.Get("/register/{inviteCode}", s.handleRegisterPage)
	s.mux.Post("/register/{inviteCode}", s.handleRegister)

	s.mux.Group(func(r chi.Router) {
		r.Use(s.requireUser)

		r.Get("/logout", s.handleLogoutPage)
		r.Post("/logout", s.handleLogout)

		r.Get("/profile", s.handleProfilePage)
		r.Post("/profile", s.handleProfileEdit)
		r.Delete("/profile", s.handleProfileDelete)
		r.Post("/profile/avatar", s.handleAvatarUpload)

		r.Get("/dashboard", s.handleDashboardPage)

		r.Group(func(r chi.Router) {
			r.Use(s.requireAdmin)

			r.Get("/organizations", s.handleOrganizationsPage)
			r.Post("/organizations", s.handleOrganizationCreate)
			r.Delete("/organizations", s.handleOrganizationDelete)
		})
	})
}

func (s *Server) render(w http.ResponseWriter, frag string, data any) {
	w.Header().Set("Vary", "HX-Request")
	if err := s.frags.ExecuteTemplate(w, frag, data); err != nil {
		log.Println("render: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) renderPage(w http.ResponseWriter, page string, data any) {
	tmpl, ok := s.pages[page]
	if !ok {
		log.Println("unknown page: ", page)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Vary", "HX-Request")
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Println("renderPage: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "index.html", nil)
}

func (s *Server) handleDashboardPage(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "dashboard.html", nil)
	// if user is admin send to admin dashboard
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}
