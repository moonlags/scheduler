-- +goose Up
CREATE TABLE employees (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    rate TEXT NOT NULL,
    profile_picture TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE employees;
