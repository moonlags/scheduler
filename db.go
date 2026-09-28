package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func openDB(path string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	goose.SetBaseFS(os.DirFS("migrations/*.sql"))
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, err
	}

	return db, db.Ping()
}
