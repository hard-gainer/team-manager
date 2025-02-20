package service

import (
	"github.com/gin-gonic/gin"
	auth "github.com/hard-gainer/task-tracker/internal/auth"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
	tmpl "github.com/hard-gainer/task-tracker/internal/template"
)

type Server struct {
	router     *gin.Engine
	store      db.Store
	authClient auth.AuthClient
}

func NewServer(store db.Store, authClient auth.AuthClient) *Server {
	server := &Server{
		store:      store,
		authClient: authClient,
	}
	router := gin.Default()

	router.SetFuncMap(tmpl.GetTemplateFuncs())
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(server.authMiddleware())
	{
		authorized.GET("/projects", server.showProjects)
		authorized.GET("/dashboard/:projectId", server.showProjectDashboard)
		authorized.GET("/dashboard", server.showDashboard)
		authorized.GET("/statistics", server.showStatistics)
		registerTaskRoutes(server, authorized)
	}

	// Auth routes
	router.GET("/login", server.showLogin)
	router.GET("/register", server.showRegister)
	router.POST("/login", server.handleLogin)
	router.POST("/register", server.handleRegister)
	router.POST("/logout", server.handleLogout)

	server.router = router
	return server
}

func registerTaskRoutes(server *Server, router *gin.RouterGroup) {
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

	router.GET("/tasks/:id/confirm", server.showTaskConfirm)
	router.GET("/tasks/:id/details", server.showTaskDetails)
	// router.GET("/tasks/create", server.showCreateTaskForm)
	// router.POST("/tasks", server.createTask)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
