package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hard-gainer/task-tracker/internal/util"
)

func (server *Server) showDashboard(ctx *gin.Context) {
	tasks, _ := server.store.ListEmployeeTasks(ctx, util.ToNullInt4(1))
	// if err != nil {
	// 	ctx.HTML(http.StatusBadRequest, "dashboatd.html", gin.H{
	// 		"tasks": tasks,
	// 	})
	// }
	ctx.HTML(http.StatusOK, "dashboard.html", gin.H{
		"tasks": tasks,
	})
}
