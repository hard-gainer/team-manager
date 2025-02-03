-- name: CreateHistory :one
INSERT INTO histories (
    task_id,
    changed_by,
    old_status,
    new_status
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: ListTaskHistory :many
SELECT * FROM histories
WHERE task_id = $1;

-- name: ListEmployeeHistory :many
SELECT * FROM histories
WHERE changed_by = $1;