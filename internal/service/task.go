package service

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	DueTo       time.Time `json:"due_to" binding:"required"`
	Status      string    `json:"status" binding:"required,oneof=ASSIGNED STARTED SUSPENDED COMPLETED"`
	Priority    string    `json:"priority" binding:"required,oneof=LOW MEDIUM HIGH CRITICAL"`
	ProjectID   int32     `json:"project_id" binding:"required"`
	AssignedTo  int32     `json:"assigned_to" binding:"required"`
}

func (server *Server) createTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTaskParams{
		Title:       req.Title,
		Description: req.Description,
		DueTo:       req.DueTo,
		Status:      req.Status,
		Priority:    req.Priority,
		ProjectID:   toNullInt4(req.ProjectID),
		AssignedTo:  toNullInt4(req.AssignedTo),
	}

	task, err := server.store.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

func (server *Server) getTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	task, err := server.store.GetTask(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func toNullInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(i),
		Valid: true,
	}
}
