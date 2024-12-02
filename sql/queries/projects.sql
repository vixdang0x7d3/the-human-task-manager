-- name: CreateProject :one
INSERT INTO projects ("id", "user_id", "title")
VALUES (@id, @user_id, @title)
RETURNING *;

-- name: ProjectByID :one
SELECT * FROM projects
WHERE id = @id;

-- name: DeleteProject :one
DELETE FROM projects
WHERE id = @id
RETURNING *;

