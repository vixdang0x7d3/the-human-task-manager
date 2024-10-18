-- +goose Up
CREATE TABLE users (
    "id"            uuid PRIMARY KEY,
    "username"      varchar(64) NOT NULL DEFAULT ''::varchar,
    "first_name"    varchar(64) NOT NULL DEFAULT ''::varchar,
    "last_name"     varchar(64) NOT NULL DEFAULT ''::varchar,
    "email"         varchar(64) NOT NULL UNIQUE,
    "password"      varchar(64) NOT NULL,
    "signup_at"     timestamptz NOT NULL,
    "last_login"    timestamptz NOT NULL
);

-- +goose Down
DROP TABLE users;
