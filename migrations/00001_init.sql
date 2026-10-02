-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    profile_picture TEXT NOT NULL,
    is_admin BOOL NOT NULL DEFAULT FALSE,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE sessions (
    token BLOB PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at INT NOT NULL,
    user_agent TEXT NOT NULL
);

CREATE TABLE organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    logo TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE roles (
    id INT PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO roles (id, name) VALUES (1, "Organization Admin");
INSERT INTO roles (id, name) VALUES (2, "Senior");
INSERT INTO roles (id, name) VALUES (3, "Employee");

-- +goose Down
DROP TABLE users;
DROP TABLE sessions;
DROP TABLE organizations;
DROP TABLE roles;
