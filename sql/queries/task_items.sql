-- name: TaskItemByID :one
SELECT * FROM task_items
WHERE user_id=$1
AND id=$2;

-- name: FindTaskItemsByUserID :many
SELECT t.*, COUNT(*) OVER()
FROM task_items t
WHERE t.user_id = @user_id
AND (
	to_tsvector(
		t.project_title
		|| ' '
		|| t.completed_by_name
		|| ' ' 
		|| t.description
		|| ' '
		|| array_to_string(tags, ' ')
	) @@ websearch_to_tsquery(sqlc.narg(q)) OR sqlc.narg(q) IS NULL
)
AND (t.state = sqlc.narg(state)::task_state OR sqlc.narg(state) IS NULL)
AND (t.priority = sqlc.narg(priority)::task_priority OR sqlc.narg(priority) IS NULL)
AND (
	   (t.deadline BETWEEN now() AND (now() + sqlc.narg(time_interval)::interval))
        OR (t.schedule BETWEEN now() AND (now() + sqlc.narg(time_interval)::interval))
        OR (t.wait BETWEEN now() AND (now() + sqlc.narg(time_interval)::interval))
	OR sqlc.narg(time_interval) IS NULL
)
ORDER BY urgency DESC
LIMIT @nlimit OFFSET @noffset;

-- name: FindTaskItemsByProjectID :many
SELECT t.*, COUNT(*) OVER()
FROM task_items t
WHERE t.project_id = @project_id
AND (
	to_tsvector(
		t.username
		|| ' '
		|| t.completed_by_name
		|| ' ' 
		|| t.description
		|| ' '
		|| array_to_string(tags, ' ')
	) @@ websearch_to_tsquery(sqlc.narg(q)) OR sqlc.narg(q) IS NULL
)
AND (t.state = sqlc.narg(state)::task_state OR sqlc.narg(state) IS NULL)
AND (t.priority = sqlc.narg(priority)::task_priority OR sqlc.narg(priority) IS NULL)
AND (
	   (t.deadline BETWEEN now() AND (now() + sqlc.narg(time_interval)::interval))
        OR (t.schedule BETWEEN now() AND (now() + sqlc.narg(time_interval)::interval))
        OR (t.wait BETWEEN now() AND (now() + sqlc.narg(time_interval)::interval))
	OR sqlc.narg(time_interval) IS NULL
)
ORDER BY urgency DESC
LIMIT @nlimit OFFSET @noffset;
