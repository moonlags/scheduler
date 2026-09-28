package main

import (
	"log"
	"net/http"
	"uuid"

	"github.com/go-chi/chi/v5"
)

type Employee struct {
	ID             string  `db:"id"`
	Name           string  `db:"name"`
	Rate           string `db:"rate"`
	ProfilePicture string  `db:"profile_picture"`
}

func (s *Server) listEmployees() ([]Employee, error) {
	var employees []Employee
	err := s.db.Select(&employees, "SELECT id, name, rate, profile_picture FROM employees")
	return employees, err
}

func (s *Server) insertEmployee(name string, rate string, profilePicture string) (Employee, error) {
	id := uuid.NewV4().String()
	if _, err := s.db.Exec(
		"INSERT INTO employees (id, name, rate, profile_picture) VALUES (?, ?, ?, ?)",
		id, name, rate, profilePicture); err != nil {
		return Employee{}, err
	}

	return Employee{ID: id, Name: name, Rate: rate, ProfilePicture: profilePicture}, nil
}

func (s *Server) deleteEmployee(id string) error {
	_, err := s.db.Exec("DELETE FROM employees WHERE id = ?", id)
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
	name := r.FormValue("name")
	rate := r.FormValue("rate")
	if name == "" {
		http.Error(w, "name is not optional", http.StatusUnprocessableEntity)
		return
	}

	c, err := s.insertEmployee(name, rate, "default")
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
