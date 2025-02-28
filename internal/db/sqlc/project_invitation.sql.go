// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: project_invitation.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const acceptProjectInvitation = `-- name: AcceptProjectInvitation :one
UPDATE project_invitations
SET accepted_at = NOW()
WHERE token = $1 AND expires_at > NOW() AND accepted_at IS NULL
RETURNING id, project_id, email, token, role, created_at, expires_at, accepted_at
`

func (q *Queries) AcceptProjectInvitation(ctx context.Context, token string) (ProjectInvitation, error) {
	row := q.db.QueryRow(ctx, acceptProjectInvitation, token)
	var i ProjectInvitation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Email,
		&i.Token,
		&i.Role,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.AcceptedAt,
	)
	return i, err
}

const createProjectInvitation = `-- name: CreateProjectInvitation :one
INSERT INTO project_invitations (
    project_id,
    email,
    token,
    role,
    expires_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, email, token, role, created_at, expires_at, accepted_at
`

type CreateProjectInvitationParams struct {
	ProjectID int64            `json:"project_id"`
	Email     string           `json:"email"`
	Token     string           `json:"token"`
	Role      string           `json:"role"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
}

func (q *Queries) CreateProjectInvitation(ctx context.Context, arg CreateProjectInvitationParams) (ProjectInvitation, error) {
	row := q.db.QueryRow(ctx, createProjectInvitation,
		arg.ProjectID,
		arg.Email,
		arg.Token,
		arg.Role,
		arg.ExpiresAt,
	)
	var i ProjectInvitation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Email,
		&i.Token,
		&i.Role,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.AcceptedAt,
	)
	return i, err
}

const getProjectInvitation = `-- name: GetProjectInvitation :one
SELECT id, project_id, email, token, role, created_at, expires_at, accepted_at FROM project_invitations
WHERE token = $1 AND expires_at > NOW() AND accepted_at IS NULL
`

func (q *Queries) GetProjectInvitation(ctx context.Context, token string) (ProjectInvitation, error) {
	row := q.db.QueryRow(ctx, getProjectInvitation, token)
	var i ProjectInvitation
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Email,
		&i.Token,
		&i.Role,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.AcceptedAt,
	)
	return i, err
}
