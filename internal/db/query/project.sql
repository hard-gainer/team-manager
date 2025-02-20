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

-- name: GetProjectWithParticipants :one
SELECT 
    p.*,
    COUNT(DISTINCT t.id) as task_count,
    COALESCE(SUM(t.time_spent), 0) as total_time_spent,
    COALESCE(
        json_agg(
            json_build_object(
                'id', e.id,
                'first_name', e.first_name,
                'last_name', e.last_name,
                'email', e.email,
                'role', pp.role
            )
        ) FILTER (WHERE e.id IS NOT NULL),
        '[]'
    ) as participants
FROM projects p
LEFT JOIN project_participants pp ON p.id = pp.project_id
LEFT JOIN employees e ON pp.user_id = e.id
LEFT JOIN tasks t ON p.id = t.project_id
WHERE p.id = $1
GROUP BY p.id;

-- name: ListUserProjectsWithParticipants :many
SELECT 
    p.*,
    COUNT(DISTINCT t.id) as task_count,
    COALESCE(SUM(t.time_spent), 0) as total_time_spent,
    COALESCE(
        json_agg(
            json_build_object(
                'id', e.id,
                'name', e.name,
                'email', e.email,
                'role', pp.role
            )
        ) FILTER (WHERE e.id IS NOT NULL),
        '[]'
    ) as participants
FROM projects p
LEFT JOIN project_participants pp ON p.id = pp.project_id
LEFT JOIN employees e ON pp.user_id = e.id
LEFT JOIN tasks t ON p.id = t.project_id
WHERE p.created_by = $1 
    OR EXISTS (
        SELECT 1 FROM project_participants 
        WHERE project_id = p.id AND user_id = $1
    )
GROUP BY p.id
ORDER BY p.start_date DESC;