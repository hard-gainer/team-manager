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

		employee, err := server.store.GetEmployee(ctx, userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or missing employee ID",
			})
			return
		}

		// Получаем ID проекта из параметров URL
		var projectID int64

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
		projectRole, err := server.getUserRoleInProject(ctx, userID, projectID)
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
				if projectRole == requiredRole {
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

		ctx.Set("app_role", employee.Role)
		ctx.Set("project_role", projectRole)
		ctx.Next()
	}
}

func checkProjectRole(ctx *gin.Context, roles ...string) bool {
	// Получаем роль пользователя из контекста
	roleValue, exists := ctx.Get("project_role")
	if !exists {
		log.Printf("Project role not found in context")
		return false
	}

	// Преобразуем в строку и проверяем
	role, ok := roleValue.(string)
	if !ok {
		log.Printf("Project role is not a string: %T", roleValue)
		return false
	}

	// Проверяем, соответствует ли роль одной из требуемых
	for _, r := range roles {
		if role == r {
			return true
		}
	}

	log.Printf("User has role '%s', but required one of: %v", role, roles)
	return false
}

func (server *Server) requireProjectRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !checkProjectRole(ctx, roles...) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions for this action",
			})
			return
		}
		ctx.Next()
	}
}

// appRoleMiddleware проверяет глобальную роль пользователя в приложении
func (server *Server) appRoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := getUserIDFromToken(ctx)

		// Получаем информацию о пользователе
		employee, err := server.store.GetEmployee(ctx, userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load user data",
			})
			return
		}

		// Получаем глобальную роль пользователя в системе
		appRole := employee.Role
		log.Printf("User %d has app role: %s", userID, appRole)

		// Проверяем, соответствует ли роль требуемым
		if len(requiredRoles) > 0 {
			hasRequiredRole := false
			for _, requiredRole := range requiredRoles {
				if appRole == requiredRole {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				log.Printf("User %d with role %s doesn't have required role (one of %v)",
					userID, appRole, requiredRoles)
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Insufficient permissions for this action",
				})
				return
			}
		}

		// Сохраняем роль в контексте для использования в обработчиках
		ctx.Set("app_role", appRole)
		ctx.Next()
	}
}

// isAppAdmin проверяет, имеет ли пользователь глобальную роль Admin
func isAppAdmin(ctx *gin.Context) bool {
	roleValue, exists := ctx.Get("app_role")
	if !exists {
		return false
	}

	role, ok := roleValue.(string)
	if !ok {
		return false
	}

	return role == AppRoleAdmin
}

// isAppManager проверяет, имеет ли пользователь глобальную роль Admin или Manager
func isAppManager(ctx *gin.Context) bool {
	roleValue, exists := ctx.Get("app_role")
	if !exists {
		return false
	}

	role, ok := roleValue.(string)
	if !ok {
		return false
	}

	return role == AppRoleAdmin || role == AppRoleManager
}

// checkAppRole проверяет, имеет ли пользователь одну из указанных глобальных ролей
func checkAppRole(ctx *gin.Context, roles ...string) bool {
	roleValue, exists := ctx.Get("app_role")
	if !exists {
		log.Printf("App role not found in context")
		return false
	}

	role, ok := roleValue.(string)
	if !ok {
		log.Printf("App role is not a string: %T", roleValue)
		return false
	}

	for _, r := range roles {
		if role == r {
			return true
		}
	}

	log.Printf("User has app role '%s', but required one of: %v", role, roles)
	return false
}
