-- +goose Up
CREATE TABLE users (
    "id"            uuid PRIMARY KEY,
    "username"      text NOT NULL DEFAULT,
    "first_name"    text NOT NULL DEFAULT,
    "last_name"     text NOT NULL DEFAULT,
    "email"         text NOT NULL UNIQUE,
    "password"      text NOT NULL,
    "signup_at"     timestamptz NOT NULL,
    "last_login"    timestamptz NOT NULL
);

-- +goose Down
DROP TABLE users;
