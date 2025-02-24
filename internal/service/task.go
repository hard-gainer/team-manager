package service

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hard-gainer/team-manager/internal/db/sqlc"
	"github.com/hard-gainer/team-manager/internal/util"
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
		ProjectID:   util.ToNullInt4(req.ProjectID),
		AssignedTo:  util.ToNullInt4(req.AssignedTo),
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

func (server *Server) getTaskTime(ctx *gin.Context) {
    id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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

    ctx.JSON(http.StatusOK, gin.H{
        "timeSpent": task.TimeSpent.Int64,
    })
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

	projectID := util.ToNullInt4(int32(id))
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

	employeeID := util.ToNullInt4(int32(id))
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

func (server *Server) updateTaskTimeSpent(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    timeParam := ctx.PostForm("time")
    if timeParam == "" {
        timeParam = ctx.Query("time")
    }

    if timeParam == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "time is required"})
        return
    }

    timeSpent, err := strconv.ParseInt(timeParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if timeSpent < 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "time must be a positive number"})
        return
    }

    arg := db.UpdateTaskTimeSpentParams{
        ID:        id,
        TimeSpent: util.ToNullInt8(timeSpent),
    }

    task, err := server.store.UpdateTaskTimeSpent(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
		"timeSpent": task.TimeSpent,
        "status":    "success",
	})
}

func (server *Server) updateTaskStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	status := ctx.PostForm("status")
	if status == "" {
		status = ctx.Query("status")
	}

	if status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	arg := db.UpdateTaskStatusParams{
		ID:     id,
		Status: status,
	}

	task, err := server.store.UpdateTaskStatus(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.HTML(http.StatusOK, "task_partial.html", task)
}

type updateTaskPriorityRequest struct {
	Priority string `json:"priority" binding:"required,oneof=LOW MEDIUM HIGH CRITICAL"`
}

func (server *Server) updateTaskPriority(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateTaskPriorityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateTaskPriorityParams{
		ID:       id,
		Priority: req.Priority,
	}

	task, err := server.store.UpdateTaskPriority(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}
