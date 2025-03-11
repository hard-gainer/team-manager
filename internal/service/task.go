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

func (server *Server) createTask(ctx *gin.Context) {
	userID := getUserIDFromToken(ctx)

	projectID, err := strconv.ParseInt(ctx.PostForm("project_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var assignedTo int32
	assigneeID := ctx.PostForm("assignee_id")
	if assigneeID != "" {
		assigneeIDInt, err := strconv.ParseInt(assigneeID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
			return
		}
		assignedTo = int32(assigneeIDInt)
	} else {
		assignedTo = userID
	}

	arg := db.CreateTaskParams{
		Title:       ctx.PostForm("title"),
		Description: ctx.PostForm("description"),
		DueTo:       util.ParseDate(ctx.PostForm("due_to")),
		Status:      "ASSIGNED",
		Priority:    ctx.PostForm("priority"),
		ProjectID:   util.ToNullInt4(int32(projectID)),
		AssignedTo:  util.ToNullInt4(assignedTo),
	}

	task, err := server.store.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "task_partial.html", task)
}

func (server *Server) showCreateTaskForm(ctx *gin.Context) {
	projectID, err := strconv.ParseInt(ctx.Query("project_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	participants, err := server.store.ListProjectParticipants(ctx, projectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load project participants"})
		return
	}

	var projectMembers []struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}

	for _, participant := range participants {
		employee, err := server.store.GetEmployee(ctx, participant.ID)
		if err != nil {
			continue
		}
		projectMembers = append(projectMembers, struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
		}{
			ID:   employee.ID,
			Name: employee.Name,
		})
	}

	ctx.HTML(http.StatusOK, "create_task_modal.html", gin.H{
		"projectID": projectID,
		"members":   projectMembers,
	})
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

func (server *Server) showTaskConfirm(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.HTML(http.StatusOK, "confirm_modal.html", gin.H{
		"ID": id,
	})
}

func (server *Server) showTaskDetails(ctx *gin.Context) {
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

	ctx.HTML(http.StatusOK, "task_details.html", task)
}
