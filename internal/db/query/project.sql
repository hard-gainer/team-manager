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

-- name: AddProjectParticipant :one
INSERT INTO project_participants (
    project_id,
    user_id,
    role
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: RemoveProjectParticipant :exec
DELETE FROM project_participants
WHERE project_id = $1 AND user_id = $2;

-- name: ListProjectParticipants :many
SELECT e.*
FROM employees e
JOIN project_participants pp ON e.id = pp.user_id
WHERE pp.project_id = $1;

-- name: UpdateParticipantRole :one
UPDATE project_participants
SET role = $3
WHERE project_id = $1 AND user_id = $2
RETURNING *;

-- name: GetProjectWithStats :many
SELECT 
    p.*,
    COUNT(DISTINCT t.id) as task_count,
    COALESCE(SUM(t.time_spent), 0)::bigint as total_time_spent
FROM projects p
LEFT JOIN tasks t ON p.id = t.project_id
GROUP BY p.id
ORDER BY p.start_date DESC;

-- name: GetProjectStats :one
SELECT 
    COUNT(DISTINCT t.id) as task_count,
    COALESCE(SUM(t.time_spent), 0)::bigint as total_time_spent
FROM tasks t
WHERE t.project_id = $1;

-- name: GetProjectParticipantRole :one
SELECT role FROM project_participants
WHERE project_id = $1 AND user_id = $2;