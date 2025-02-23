package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	auth "github.com/hard-gainer/task-tracker/internal/auth"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
)

func (server *Server) syncUser(ctx context.Context, userID int32) error {
	userResp, err := server.authClient.GetUser(ctx, &auth.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return fmt.Errorf("failed to get user from auth service: %v", err)
	}

	_, err = server.store.CreateEmployee(ctx, db.CreateEmployeeParams{
		ID:    userResp.Id,
		Name:  userResp.Name,
		Email: userResp.Email,
		Role:  "NOT SET",
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			_, err = server.store.UpdateEmployeeName(ctx, db.UpdateEmployeeNameParams{
				ID:   userResp.Id,
				Name: userResp.Name,
			})
			if err != nil {
				return fmt.Errorf("failed to update employee: %v", err)
			}
		} else {
			return fmt.Errorf("failed to create employee: %v", err)
		}
	}

	return nil
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

	resp, err := server.authClient.Register(context.Background(), authReq)
	if err != nil {
		ctx.HTML(http.StatusOK, "register.html", gin.H{
			"error": "Failed to register user",
		})
		return
	}

	if err := server.syncUser(ctx, int32(resp.UserId)); err != nil {
		log.Printf("Failed to sync new user: %v", err)
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
		fmt.Printf("Login error: %v\n", err)
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
