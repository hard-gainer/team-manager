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

func (server *Server) listTasks(ctx *gin.Context) {
    tasks, err := server.store.ListTasks(ctx)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, tasks)
}

func (server *Server) listProjectTasks(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    projectID := toNullInt4(int32(id))
    tasks, err := server.store.ListProjectTasks(ctx, projectID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, tasks)
}

func (server *Server) listEmployeeTasks(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    employeeID := toNullInt4(int32(id))
    tasks, err := server.store.ListEmployeeTasks(ctx, employeeID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, tasks)
}

type updateTaskTitleRequest struct {
    Title string `json:"title" binding:"required"`
}

func (server *Server) updateTaskTitle(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    var req updateTaskTitleRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.UpdateTaskTitleParams{
        ID:    id,
        Title: req.Title,
    }

    task, err := server.store.UpdateTaskTitle(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, task)
}

type updateTaskDescriptionRequest struct {
    Description string `json:"description" binding:"required"`
}

func (server *Server) updateTaskDescription(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    var req updateTaskDescriptionRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.UpdateTaskDescriptionParams{
        ID:          id,
        Description: req.Description,
    }

    task, err := server.store.UpdateTaskDescription(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, task)
}

type updateTaskDeadlineRequest struct {
    DueTo time.Time `json:"due_to" binding:"required"`
}

func (server *Server) updateTaskDeadline(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    var req updateTaskDeadlineRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.UpdateTaskDeadlineParams{
        ID:    id,
        DueTo: req.DueTo,
    }

    task, err := server.store.UpdateTaskDeadline(ctx, arg)
    if err != nil {
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
