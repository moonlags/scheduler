package main

import (
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	Email          string `db:"email"`
	Password       string `db:"password"`
	ProfilePicture string `db:"profile_picture"`
	IsAdmin        bool   `db:"is_admin"`
}

func (s *Server) insertUser(name string, email string, password string, profilePicture string, isAdmin bool) (User, error) {
	id := ulid.Make()
	if _, err := s.db.Exec(
		"INSERT INTO users (id, name, email, password, profile_picture, is_admin) VALUES (?, ?, ?, ?, ?, ?)",
		id.String(), name, email, password, profilePicture, isAdmin); err != nil {
		return User{}, err
	}

	return User{ID: id.String(), Name: name, Email: email, Password: password, ProfilePicture: profilePicture, IsAdmin: isAdmin}, nil
}

func (s *Server) getUserByEmail(email string) (User, error) {
	var user User
	if err := s.db.Get(&user, "SELECT FROM users WHERE email = ?", email); err != nil{
		return User{}, err
	}
	return user, nil
}
