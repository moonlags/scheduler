package store

type Employee struct {
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Rate   string `db:"rate"`
}

const employeeQuery = `
	SELECT e.user_id, u.name, e.rate
	FROM employees e
	JOIN users u ON u.id = e.user_id`

func (s *Store) Employees() ([]Employee, error) {
	var employees []Employee
	err := s.db.Select(&employees, employeeQuery+" ORDER BY u.name")
	return employees, err
}

func (s *Store) EmployeeByID(userID string) (Employee, error) {
	var e Employee
	err := s.db.Get(&e, employeeQuery+" WHERE e.user_id = ?", userID)
	return e, err
}

func (s *Store) InsertEmployee(userID, rate string) (Employee, error) {
	if _, err := s.db.Exec(
		"INSERT INTO employees (user_id, rate) VALUES (?, ?)",
		userID, rate); err != nil {
		return Employee{}, err
	}
	return s.EmployeeByID(userID)
}

func (s *Store) DeleteEmployee(userID string) error {
	_, err := s.db.Exec("DELETE FROM employees WHERE user_id = ?", userID)
	return err
}
