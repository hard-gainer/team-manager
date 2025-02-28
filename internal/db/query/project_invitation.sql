-- name: CreateProjectInvitation :one
INSERT INTO project_invitations (
    project_id,
    email,
    token,
    role,
    expires_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProjectInvitation :one
SELECT * FROM project_invitations
WHERE token = $1 AND expires_at > NOW() AND accepted_at IS NULL;

-- name: AcceptProjectInvitation :one
UPDATE project_invitations
SET accepted_at = NOW()
WHERE token = $1 AND expires_at > NOW() AND accepted_at IS NULL
RETURNING *;