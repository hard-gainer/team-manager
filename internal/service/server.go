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
		// Общие маршруты для всех авторизованных
		authorized.GET("/projects", server.showProjects)

		// Перенаправление с /projects/:id на dashboard
		authorized.GET("/projects/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			ctx.Redirect(http.StatusFound, "/dashboard/"+id)
		})

		// Dashboard проекта с учетом роли
		authorized.GET("/dashboard/:projectId", server.projectRoleMiddleware(ProjectRoleMember, ProjectRoleManager, ProjectRoleOwner), server.showProjectDashboard)

		// Маршруты только для менеджеров и владельцев (управление проектом)
		projectManagerRoutes := authorized.Group("/")
		projectManagerRoutes.Use(server.requireManagementRights())
		{
			// Создание проектов
			projectManagerRoutes.GET("/projects/create", server.showCreateProjectForm)
			projectManagerRoutes.POST("/projects", server.createProject)

			// Статистика
			projectManagerRoutes.GET("/statistics", server.showStatistics)
		}

		// Маршруты для задач (просмотр, обновление)
		taskRoutes := authorized.Group("/tasks")
		{
			taskRoutes.GET("", server.listTasks)
			taskRoutes.GET("/:id", server.getTask)
			taskRoutes.GET("/:id/time", server.getTaskTime)
			taskRoutes.GET("/:id/confirm", server.showTaskConfirm)
			taskRoutes.GET("/:id/details", server.showTaskDetails)

			// Обновление задач (любой участник может обновлять свои задачи)
			taskRoutes.PATCH("/:id/time", server.updateTaskTimeSpent)
			taskRoutes.PATCH("/:id/status", server.updateTaskStatus)
		}

		// Маршруты для операций с задачами, доступные только менеджерам/владельцам
		managerTaskRoutes := authorized.Group("/tasks")
		managerTaskRoutes.Use(server.requireManagementRights())
		{
			managerTaskRoutes.GET("/create", server.showCreateTaskForm)
			managerTaskRoutes.POST("", server.createTask)

			// Редактирование задач доступно только менеджерам
			managerTaskRoutes.PATCH("/:id/title", server.updateTaskTitle)
			managerTaskRoutes.PATCH("/:id/description", server.updateTaskDescription)
			managerTaskRoutes.PATCH("/:id/deadline", server.updateTaskDeadline)
			managerTaskRoutes.PATCH("/:id/priority", server.updateTaskPriority)
		}

		// Маршруты для управления проектом
		projectManagementRoutes := authorized.Group("/projects/:id")
		projectManagementRoutes.Use(server.projectRoleMiddleware(ProjectRoleManager, ProjectRoleOwner))
		{
			projectManagementRoutes.GET("/invite", server.showInviteMemberForm)
			projectManagementRoutes.POST("/invite", server.inviteMember)
		}

		// Приглашения в проект и работа с сотрудниками
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
