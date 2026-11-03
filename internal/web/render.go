package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
)

// Templates and static assets are compiled into the binary.
//
//go:embed templates static
var assets embed.FS

func (s *Server) loadTemplates() error {
	pages, err := fs.Glob(assets, "templates/pages/*.html")
	if err != nil {
		return fmt.Errorf("list page templates: %w", err)
	}

	s.pages = make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		tmpl, err := template.ParseFS(assets, "templates/layout.html", "templates/fragments/*.html", page)
		if err != nil {
			return fmt.Errorf("parse page %s: %w", page, err)
		}
		s.pages[path.Base(page)] = tmpl
	}

	frags, err := template.ParseFS(assets, "templates/fragments/*.html")
	if err != nil {
		return fmt.Errorf("parse fragments: %w", err)
	}
	s.frags = frags

	return nil
}

// render writes a single template fragment (used for HTMX partial responses).
func (s *Server) render(w http.ResponseWriter, frag string, data any) {
	w.Header().Set("Vary", "HX-Request")
	if err := s.frags.ExecuteTemplate(w, frag, data); err != nil {
		log.Println("render:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

// renderPage writes a full page wrapped in the layout.
func (s *Server) renderPage(w http.ResponseWriter, page string, data any) {
	tmpl, ok := s.pages[page]
	if !ok {
		log.Println("unknown page:", page)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Vary", "HX-Request")
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Println("renderPage:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
