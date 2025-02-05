package service

import (
	"github.com/gin-gonic/gin"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
)

type Server struct {
	router *gin.Engine
	store db.Store
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/tasks/:id", server.getTask)
    router.GET("/tasks", server.listTasks)
    router.GET("/projects/:id/tasks", server.listProjectTasks)
    router.GET("/employees/:id/tasks", server.listEmployeeTasks)
	router.POST("/tasks", server.createTask)
    router.PATCH("/tasks/:id/title", server.updateTaskTitle)
   	router.PATCH("/tasks/:id/description", server.updateTaskDescription)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
    return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}