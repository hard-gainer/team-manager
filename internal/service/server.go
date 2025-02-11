package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
	tmpl "github.com/hard-gainer/task-tracker/internal/template"
	"github.com/hard-gainer/task-tracker/internal/util"
)

type Server struct {
	router *gin.Engine
	store  db.Store
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.SetFuncMap(tmpl.GetTemplateFuncs())
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	registerTaskRoutes(server, router)
	router.GET("/dashboard", server.showDashboard)
	router.GET("/statistics", server.showStatistics)

	server.router = router
	return server
}

func registerTaskRoutes(server *Server, router *gin.Engine) {
	router.GET("/tasks/:id", server.getTask)
	router.GET("/tasks", server.listTasks)
	router.GET("/projects/:id/tasks", server.listProjectTasks)
	router.GET("/employees/:id/tasks", server.listEmployeeTasks)
	router.POST("/tasks", server.createTask)
	router.PATCH("/tasks/:id/title", server.updateTaskTitle)
	router.PATCH("/tasks/:id/description", server.updateTaskDescription)
	router.PATCH("/tasks/:id/deadline", server.updateTaskDeadline)
	router.PATCH("/tasks/:id/time", server.updateTaskTimeSpent)
	router.PATCH("/tasks/:id/status", server.updateTaskStatus)
	router.PATCH("/tasks/:id/priority", server.updateTaskPriority)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

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

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
