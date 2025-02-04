-- name: CreateProject :one
INSERT INTO projects (
    title,
    description,
    start_date,
    end_date,
    created_by
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = $1;

-- name: ListProjects :many
SELECT * FROM projects;

-- name: UpdateProject :one
UPDATE projects
    SET title = $2,
    end_date = $3
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;