package web

import "net/http"

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "index.html", nil)
}

func (s *Server) handleDashboardPage(w http.ResponseWriter, r *http.Request) {
	// TODO: render an admin dashboard for administrators.
	s.renderPage(w, "dashboard.html", userFrom(r.Context()))
}

// redirect sends HTMX requests to a new page via a response header and other
// requests via a regular redirect.
func redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}
