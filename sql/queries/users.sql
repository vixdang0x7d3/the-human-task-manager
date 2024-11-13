-- name: CreateUser :one
INSERT INTO users(id, username, first_name, last_name, email, password, signup_at, last_login)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ByID :one
SELECT * FROM users
WHERE id=$1;

-- name: ByEmail :one
SELECT * FROM users
WHERE email=$1;
