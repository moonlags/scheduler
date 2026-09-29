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
)

type Server struct {
	mux   *chi.Mux
	db    *sqlx.DB
	pages map[string]*template.Template
	frags *template.Template
}

func newServer(db *sqlx.DB) (*Server, error) {
	server := &Server{
		mux: chi.NewRouter(),
		db:  db,
	}
	server.loadTemplates(os.DirFS("."))
	server.routes()

	return server, nil
}

func (s *Server) run(addr string) error {
	log.Printf("running on address %s\n", addr)

	return http.ListenAndServe(addr, s.mux)
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
	s.mux.Use(func(next http.Handler) http.Handler { return cop.Handler(next)})

	s.mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	s.mux.Get("/", s.handleIndex)

	s.mux.Get("/login", s.handleLoginPage)
	s.mux.Post("/login", s.handleLogin)
	s.mux.Get("/register", s.handleRegisterPage)
	s.mux.Post("/register", s.handleRegister)

	s.mux.Get("/employees", s.handleEmployeesPage)
	s.mux.Post("/employees", s.handleEmployeeCreate)
	s.mux.Delete("/employees/{id}", s.handleEmployeeDelete)
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

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "index.html", nil)
}
