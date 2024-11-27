-- name: TaskItemByID :one
SELECT * FROM task_items
WHERE id=$1;
