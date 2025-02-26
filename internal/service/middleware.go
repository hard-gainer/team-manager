package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hard-gainer/team-manager/internal/auth"
)

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
