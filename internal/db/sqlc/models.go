// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Employee struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type History struct {
	ID        int64       `json:"id"`
	TaskID    pgtype.Int4 `json:"task_id"`
	ChangedBy pgtype.Int4 `json:"changed_by"`
	ChangeAt  time.Time   `json:"change_at"`
	OldStatus string      `json:"old_status"`
	NewStatus string      `json:"new_status"`
}

type Project struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	StartDate   time.Time   `json:"start_date"`
	EndDate     time.Time   `json:"end_date"`
	CreatedBy   pgtype.Int4 `json:"created_by"`
}

type Task struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	DueTo       time.Time   `json:"due_to"`
	TimeSpent   pgtype.Int8 `json:"time_spent"`
	Status      string      `json:"status"`
	Priority    string      `json:"priority"`
	ProjectID   pgtype.Int4 `json:"project_id"`
	AssignedTo  pgtype.Int4 `json:"assigned_to"`
}
