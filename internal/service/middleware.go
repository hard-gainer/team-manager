package service

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func (server *Server) authMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token, err := ctx.Cookie("auth_token")
        if err != nil || token == "" {
            ctx.Redirect(http.StatusSeeOther, "/login")
            ctx.Abort()
            return
        }
        
        ctx.Next()
    }
}