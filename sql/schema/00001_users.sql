-- +goose UP
CREATE TABLE users (
    "id"            uuid PRIMARY KEY,
    "first_name"    varchar(64) DEFAULT ''::varchar,
    "last_name"     varchar(64) DEFAULT ''::varchar,
    "email"         varchar(64) NOT NULL,
    "password"      varchar(64) NOT NULL,
    "signup_at"     timestamptz DEFAULT now(),
    "last_login"    timestamptz DEFAULT now()
);

-- +goose DOWN
DROP TABLE users;