package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleEmployeesPage(w http.ResponseWriter, r *http.Request) {
	employees, err := s.store.Employees()
	if err != nil {
		log.Println("handleEmployeesPage:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	s.renderPage(w, "employees.html", employees)
}

func (s *Server) handleEmployeeCreate(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	rate := r.FormValue("rate")
	if userID == "" {
		http.Error(w, "user_id is not optional", http.StatusUnprocessableEntity)
		return
	}

	e, err := s.store.InsertEmployee(userID, rate)
	if err != nil {
		log.Println("handleEmployeeCreate:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	s.render(w, "employee_row.html", e)
}

func (s *Server) handleEmployeeDelete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if err := s.store.DeleteEmployee(userID); err != nil {
		log.Println("handleEmployeeDelete:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
