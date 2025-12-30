-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
    (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
    (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

ALTER TABLE users
    ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS apps;
-- +goose StatementEnd
