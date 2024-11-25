-- name: CreateTask :one
INSERT INTO tasks("id", "user_id", "project_id", "description", "priority", "status", "deadline", "schedule", "wait", "create", "end")
VALUES ($1, $2, $3,  $4, @priority::task_priority, @status::task_status, $5, $6, $7, $8, $9)
RETURNING *;

