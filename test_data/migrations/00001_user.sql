-- +goose Up
CREATE TABLE users (
    "id"            uuid PRIMARY KEY,
    "username"      text NOT NULL,
    "first_name"    text NOT NULL,
    "last_name"     text NOT NULL,
    "email"         text NOT NULL UNIQUE,
    "password"      text NOT NULL,
    "signup_at"     timestamptz NOT NULL,
    "last_login"    timestamptz NOT NULL
);

-- +goose Down
DROP TABLE users;
