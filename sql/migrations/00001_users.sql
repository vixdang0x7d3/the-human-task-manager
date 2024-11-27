-- +goose Up
CREATE TABLE users (
    "id"            uuid PRIMARY KEY,
    "username"      text NOT NULL DEFAULT ''::varchar,
    "first_name"    text NOT NULL DEFAULT ''::varchar,
    "last_name"     text NOT NULL DEFAULT ''::varchar,
    "email"         text NOT NULL UNIQUE,
    "password"      text NOT NULL,
    "signup_at"     timestamptz NOT NULL,
    "last_login"    timestamptz NOT NULL
);

-- +goose Down
DROP TABLE users;
