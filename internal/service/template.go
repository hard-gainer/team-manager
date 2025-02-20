package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	auth "github.com/hard-gainer/task-tracker/internal/auth"
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

func (server *Server) handleRegister(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	name := ctx.PostForm("name")

	if email == "" || password == "" || name == "" {
		ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "All fields are required",
		})
		return
	}

	authReq := &auth.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
		IsAdmin:  false,
	}

	_, err := server.authClient.Register(context.Background(), authReq)
	if err != nil {
		ctx.HTML(http.StatusOK, "register.html", gin.H{
			"error": "Failed to register user",
		})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/login")
}

func (server *Server) handleLogin(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if email == "" || password == "" {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Email and password are required",
		})
		return
	}

	authReq := &auth.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    1,
	}

	resp, err := server.authClient.Login(context.Background(), authReq)
	if err != nil {
		fmt.Printf("Login error: %v\n", err) // Добавляем детальное логирование
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	ctx.SetCookie("auth_token", resp.Token, 3600*24, "/", "", false, true)

	validateResp, err := server.authClient.ValidateToken(context.Background(), &auth.ValidateTokenRequest{
		Token: resp.Token,
	})
	if err != nil || !validateResp.IsValid {
		fmt.Printf("Token validation error: %v\n", err)
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Authentication failed",
		})
		return
	}

	fmt.Printf("Login successful, redirecting with token: %s\n", resp.Token)
	ctx.Redirect(http.StatusSeeOther, "/projects")
}

func (server *Server) handleLogout(ctx *gin.Context) {
	token, err := ctx.Cookie("auth_token")
	if err == nil {
		_, _ = server.authClient.Logout(context.Background(), &auth.LogoutRequest{Token: token})
	}

	ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
	ctx.Redirect(http.StatusSeeOther, "/login")
}

func getUserIDFromToken(ctx *gin.Context) int32 {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}

	if id, ok := userID.(int32); ok {
		return id
	}
	return 0
}
