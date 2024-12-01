-- name: CreateProject :one
INSERT INTO projects ("id", "user_id", "title")
VALUES (@id, @user_id, @title)
RETURNING *;
