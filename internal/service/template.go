package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hard-gainer/team-manager/internal/util"
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

// func (server *Server) showCreateTaskForm(ctx *gin.Context) {
// 	ctx.HTML(http.StatusOK, "create_task_modal.html", nil)
// }

func (server *Server) showLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}

func (server *Server) showRegister(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}
