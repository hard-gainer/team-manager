package service

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hard-gainer/task-tracker/internal/util"
)

func (server *Server) showStatistics(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "statistics.html", gin.H{
		"active": "statistics",
	})
}

func (server *Server) showDashboard(ctx *gin.Context) {
	tasks, _ := server.store.ListEmployeeTasks(ctx, util.ToNullInt4(1))
	ctx.HTML(http.StatusOK, "dashboard.html", gin.H{
		"tasks":  tasks,
		"active": "dashboard",
	})
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

// func (server *Server) showCreateTaskForm(ctx *gin.Context) {
// 	ctx.HTML(http.StatusOK, "create_task_modal.html", nil)
// }

func (server *Server) showProjects(ctx *gin.Context) {
	userID := getUserIDFromToken(ctx)

	projects, err := server.store.ListUserProjectsWithParticipants(ctx, util.ToNullInt4(userID))
	if err != nil {
		log.Printf("Error loading projects: %v", err)

		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load projects",
		})
		return
	}

	ctx.HTML(http.StatusOK, "projects.html", gin.H{
		"projects": projects,
		"active":   "projects",
	})
}

func (server *Server) showProjectDashboard(ctx *gin.Context) {
	projectID, err := strconv.ParseInt(ctx.Param("projectId"), 10, 32)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/projects")
		return
	}

	tasks, err := server.store.ListProjectTasks(ctx, util.ToNullInt4(int32(projectID)))
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/projects")
		return
	}

	ctx.HTML(http.StatusOK, "dashboard.html", gin.H{
		"tasks":     tasks,
		"active":    "dashboard",
		"projectID": projectID,
	})
}

// Auth handlers
func (server *Server) showLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}

func (server *Server) showRegister(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}
