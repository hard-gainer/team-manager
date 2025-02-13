package service

import (
	"github.com/gin-gonic/gin"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
	tmpl "github.com/hard-gainer/task-tracker/internal/template"
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
    router.GET("/tasks/:id/time", server.getTaskTime)
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

	// html routes	
	router.GET("/tasks/:id/confirm", server.showTaskConfirm)
	router.GET("/tasks/:id/details", server.showTaskDetails)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
