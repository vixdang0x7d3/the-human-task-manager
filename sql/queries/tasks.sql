-- name: CreateTask :one
INSERT INTO tasks(
	"id", 
	"user_id", 
	"project_id", 
	"description", 
	"priority", 
	"state", 
	"deadline", 
	"schedule", 
	"wait", 
	"create", 
	"end", 
	"tags"
)
VALUES ($1, $2, $3,  $4, @priority::task_priority, @state::task_state, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: TaskByID :one
SELECT * FROM tasks
WHERE id = $1;


-- name: UpdateTask :one
UPDATE tasks
SET 
	description = @description,
	priority = @priority::task_priority,
	state = @state::task_state,
	deadline = @deadline,
	schedule = @schedule,
	wait = @wait,
	tags = @tags
WHERE id = @id
RETURNING *;

-- name: SetTaskProject :one
UPDATE tasks
SET
	project_id = sqlc.narg(project_id)
WHERE id = @id
RETURNING *;

-- name: CompleteTask :one
UPDATE tasks
SET
	state = 'completed'::task_state,
	completed_by = @user_id,
	"end" = @end_timestamp
WHERE id = @id
RETURNING *;

-- name: DeleteTask :one
UPDATE tasks
SET
	state = 'deleted'::task_state
WHERE id = @id
RETURNING *;

-- name: StartTasks :many
UPDATE tasks
SET
	state = 'started'::task_state
WHERE 
	state = 'waiting'::task_state
AND	wait <= now()
RETURNING *;

-- name: StartTask :one
UPDATE tasks
SET
	state = 'started'::task_state
WHERE id= @id
RETURNING *;



