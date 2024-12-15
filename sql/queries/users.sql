-- name: CreateUser :one
INSERT INTO users("id", "username", "first_name", "last_name", "email", "password", "signup_at", "last_login")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
	username=$2,
	first_name=$3,
	last_name=$4,
	email=$5,
	password=$6
WHERE id=$1
RETURNING *;

-- name: UserByID :one
SELECT * FROM users
WHERE id=$1;

-- name: UserByEmail :one
SELECT * FROM users
WHERE email=$1;

-- name: DeleteUser :one
DELETE FROM users
WHERE id=$1
RETURNING *;