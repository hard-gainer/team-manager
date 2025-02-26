package types

import (
	db "github.com/hard-gainer/team-manager/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectWithStats struct {
	db.Project
	TaskCount int64       `json:"task_count"`
	TotalTimeSpent pgtype.Int8 `json:"time_spent"`
}
