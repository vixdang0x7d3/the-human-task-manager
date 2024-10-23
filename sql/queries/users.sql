-- name: CreateUser :one
INSERT INTO users(id, username, first_name, last_name, email, password, signup_at, last_login)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id=$1;

