package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/hard-gainer/team-manager/internal/auth"
	"github.com/hard-gainer/team-manager/internal/config"
	db "github.com/hard-gainer/team-manager/internal/db/sqlc"
	"github.com/hard-gainer/team-manager/internal/mail"
	tmpl "github.com/hard-gainer/team-manager/internal/template"
)

type Server struct {
	router     *gin.Engine
	store      db.Store
	authClient auth.AuthClient
	config     *config.Config
	mailer     *mail.Mailer
}

func NewServer(
	cfg *config.Config,
	store db.Store,
	authClient auth.AuthClient,
	mailer *mail.Mailer,
) *Server {
	server := &Server{
		store:      store,
		authClient: authClient,
		config:     cfg,
		mailer:     mailer,
	}
	router := gin.Default()

	router.SetFuncMap(tmpl.GetTemplateFuncs())
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Auth routes
	router.GET("/login", server.showLogin)
	router.GET("/register", server.showRegister)
	router.POST("/login", server.handleLogin)
	router.POST("/register", server.handleRegister)
	router.POST("/logout", server.handleLogout)

	// Protected routes - требуется аутентификация
	authorized := router.Group("/")
	authorized.Use(server.authMiddleware())
	{
		authorized.GET("/projects", server.showProjects)
		authorized.GET("/projects/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			ctx.Redirect(http.StatusFound, "/dashboard/"+id)
		})

		authorized.GET("/dashboard/:projectId", server.projectRoleMiddleware(ProjectRoleMember, ProjectRoleManager, ProjectRoleOwner), server.showProjectDashboard)

		// project routes, permission: admin, manager
		appManagerRoutes := authorized.Group("/")
		appManagerRoutes.Use(server.appRoleMiddleware(AppRoleAdmin, AppRoleManager))
		{
			appManagerRoutes.GET("projects/create", server.showCreateProjectForm)
			appManagerRoutes.POST("projects", server.createProject)

			appManagerRoutes.GET("statistics", server.showStatistics)
		}

		// tasks routs, permission: all
		taskRoutes := authorized.Group("/tasks")
		{
			taskRoutes.GET("", server.listTasks)
			taskRoutes.GET("/:id", server.getTask)
			taskRoutes.GET("/:id/time", server.getTaskTime)
			taskRoutes.GET("/:id/confirm", server.showTaskConfirm)
			taskRoutes.GET("/:id/details", server.showTaskDetails)

			taskRoutes.PATCH("/:id/time", server.updateTaskTimeSpent)
			taskRoutes.PATCH("/:id/status", server.updateTaskStatus)
		}

		// tasks routs, permission: owner, manager
		managerTaskRoutes := authorized.Group("")
		managerTaskRoutes.Use(server.projectRoleMiddleware(ProjectRoleOwner, ProjectRoleManager))
		{
			managerTaskRoutes.GET("/tasks/create", server.showCreateTaskForm)
			managerTaskRoutes.POST("/tasks", server.createTask)

			managerTaskRoutes.PATCH("/tasks/:id/title", server.updateTaskTitle)
			managerTaskRoutes.PATCH("/tasks/:id/description", server.updateTaskDescription)
			managerTaskRoutes.PATCH("/tasks/:id/deadline", server.updateTaskDeadline)
			managerTaskRoutes.PATCH("/tasks/:id/priority", server.updateTaskPriority)

			managerTaskRoutes.GET("/projects/:id/invite", server.showInviteMemberForm)
			managerTaskRoutes.POST("/projects/:id/invite", server.inviteMember)
		}

		authorized.GET("/projects/join/:token", server.handleProjectInvitation)
		authorized.GET("/employees/:id/tasks", server.listEmployeeTasks)
	}

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
