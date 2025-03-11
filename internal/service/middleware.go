package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hard-gainer/team-manager/internal/auth"
)

// authMiddleware проверяет аутентификацию пользователя
func (server *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("auth_token")
		if err != nil || token == "" {
			fmt.Printf("No auth token found: %v\n", err)
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		resp, err := server.authClient.ValidateToken(context.Background(), &auth.ValidateTokenRequest{
			Token: token,
		})
		if err != nil || !resp.IsValid {
			fmt.Printf("Token validation failed: %v\n", err)
			ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		if err := server.syncUser(ctx, resp.UserId); err != nil {
			log.Printf("Failed to sync user: %v", err)
		}

		ctx.Set("user_id", resp.UserId)
		ctx.Next()
	}
}

// projectRoleMiddleware проверяет роль пользователя в проекте
func (server *Server) projectRoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := getUserIDFromToken(ctx)

		// Получаем ID проекта из параметров URL
		var projectID int64
		var err error

		// Пробуем получить из :id
		if idParam := ctx.Param("id"); idParam != "" {
			projectID, err = strconv.ParseInt(idParam, 10, 64)
		} else if idParam := ctx.Param("projectId"); idParam != "" {
			// Или из :projectId
			projectID, err = strconv.ParseInt(idParam, 10, 64)
		} else {
			// Или из параметров запроса (для task_create и подобных)
			projectIDStr := ctx.Query("project_id")
			if projectIDStr != "" {
				projectID, err = strconv.ParseInt(projectIDStr, 10, 64)
			}
		}

		if err != nil || projectID == 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or missing project ID",
			})
			return
		}

		// Получаем роль пользователя в проекте
		role, err := server.getUserRoleInProject(ctx, userID, projectID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You don't have access to this project",
			})
			return
		}

		// Проверяем, соответствует ли роль требуемым
		if len(requiredRoles) > 0 {
			hasRequiredRole := false
			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Insufficient permissions for this action",
				})
				return
			}
		}

		// Сохраняем роль в контексте для использования в обработчиках
		ctx.Set("project_role", role)
		ctx.Next()
	}
}

// requireRoleMiddleware проверяет роль для операций, не привязанных к проекту
func (server *Server) requireManagementRights() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := getUserIDFromToken(ctx)

		// Получаем информацию о пользователе
		employee, err := server.store.GetEmployee(ctx, userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user information",
			})
			return
		}

		// Проверяем, имеет ли пользователь права manager или admin
		if employee.Role != "admin" && employee.Role != "manager" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You don't have permission to perform this action",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
