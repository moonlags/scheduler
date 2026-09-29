package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Employee struct {
	UserID string `db:"user_id"`
	Rate   string `db:"rate"`
}

func (s *Server) listEmployees() ([]Employee, error) {
	var employees []Employee
	err := s.db.Select(&employees, "SELECT user_id,  rate FROM employees")
	return employees, err
}

func (s *Server) insertEmployee(user_id string, rate string) (Employee, error) {
	if _, err := s.db.Exec(
		"INSERT INTO employees (user_id, rate) VALUES (?, ?)",
		user_id, rate); err != nil {
		return Employee{}, err
	}

	return Employee{UserID: user_id, Rate: rate}, nil
}

func (s *Server) deleteEmployee(id string) error {
	_, err := s.db.Exec("DELETE FROM employees WHERE user_id = ?", id)
	return err
}

func (s *Server) handleEmployeesPage(w http.ResponseWriter, r *http.Request) {
	employees, err := s.listEmployees()
	if err != nil {
		log.Println("handleEmployeesPage: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	s.renderPage(w, "employees.html", employees)
}

func (s *Server) handleEmployeeCreate(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	rate := r.FormValue("rate")
	if userID == "" {
		http.Error(w, "name is not optional", http.StatusUnprocessableEntity)
		return
	}

	c, err := s.insertEmployee(userID, rate)
	if err != nil {
		log.Println("handleEmployeesCreate: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	s.render(w, "employee_row.html", c)
}

func (s *Server) handleEmployeeDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.deleteEmployee(id); err != nil {
		log.Println("handleEmployeesDelete: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
