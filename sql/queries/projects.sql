-- name: CreateProject :one
INSERT INTO projects ("id", "user_id", "title")
VALUES (@id, @user_id, @title)
RETURNING *;

-- name: ProjectByID :one
SELECT * FROM projects
WHERE id = @id;

-- name: ProjectsByUserID :many
SELECT p.*, COUNT(*) OVER()  FROM projects p
WHERE user_id = @user_id
LIMIT @nlimit OFFSET @noffset;

-- name: DeleteProject :one
DELETE FROM projects
WHERE id = @id
RETURNING *;


