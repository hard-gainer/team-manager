// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	AddProjectParticipant(ctx context.Context, arg AddProjectParticipantParams) (ProjectParticipant, error)
	CreateEmployee(ctx context.Context, arg CreateEmployeeParams) (Employee, error)
	CreateHistory(ctx context.Context, arg CreateHistoryParams) (History, error)
	CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	DeleteEmployee(ctx context.Context, id int32) error
	DeleteProject(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
	GetEmployee(ctx context.Context, id int32) (Employee, error)
	GetProject(ctx context.Context, id int64) (Project, error)
	GetTask(ctx context.Context, id int64) (Task, error)
	ListEmployeeHistory(ctx context.Context, changedBy pgtype.Int4) ([]History, error)
	ListEmployeeTasks(ctx context.Context, assignedTo pgtype.Int4) ([]Task, error)
	ListEmployees(ctx context.Context) ([]Employee, error)
	ListProjectParticipants(ctx context.Context, projectID int64) ([]Employee, error)
	ListProjectTasks(ctx context.Context, projectID pgtype.Int4) ([]Task, error)
	ListProjects(ctx context.Context) ([]Project, error)
	ListTaskHistory(ctx context.Context, taskID pgtype.Int4) ([]History, error)
	ListTasks(ctx context.Context) ([]Task, error)
	RemoveProjectParticipant(ctx context.Context, arg RemoveProjectParticipantParams) error
	UpdateEmployee(ctx context.Context, arg UpdateEmployeeParams) (Employee, error)
	UpdateEmployeeEmail(ctx context.Context, arg UpdateEmployeeEmailParams) (Employee, error)
	UpdateEmployeeName(ctx context.Context, arg UpdateEmployeeNameParams) (Employee, error)
	UpdateEmployeeRole(ctx context.Context, arg UpdateEmployeeRoleParams) (Employee, error)
	UpdateParticipantRole(ctx context.Context, arg UpdateParticipantRoleParams) (ProjectParticipant, error)
	UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error)
	UpdateTaskDeadline(ctx context.Context, arg UpdateTaskDeadlineParams) (Task, error)
	UpdateTaskDescription(ctx context.Context, arg UpdateTaskDescriptionParams) (Task, error)
	UpdateTaskPriority(ctx context.Context, arg UpdateTaskPriorityParams) (Task, error)
	UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) (Task, error)
	UpdateTaskTimeSpent(ctx context.Context, arg UpdateTaskTimeSpentParams) (Task, error)
	UpdateTaskTitle(ctx context.Context, arg UpdateTaskTitleParams) (Task, error)
}

var _ Querier = (*Queries)(nil)
