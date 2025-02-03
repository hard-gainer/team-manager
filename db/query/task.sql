-- name: CreateTask :one
INSERT INTO tasks (
    title,
    description,
    due_to,
    status,
    priority,
    project_id,
    assigned_to
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks;

-- name: ListProjectTasks :many
SELECT * FROM tasks
WHERE project_id = $1;

-- name: ListEmployeeTasks :many
SELECT * FROM tasks
WHERE assigned_to = $1;

-- name: UpdateTask :one
UPDATE tasks
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
