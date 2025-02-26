package service

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/hard-gainer/team-manager/internal/db/sqlc"
	"github.com/hard-gainer/team-manager/internal/db/types"
	"github.com/hard-gainer/team-manager/internal/util"
)

func (server *Server) showProjects(ctx *gin.Context) {
	userID := getUserIDFromToken(ctx)

	projectsWithStats, err := server.store.GetProjectWithStats(ctx)
	if err != nil {
		log.Printf("Error loading projects: %v", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load projects",
		})
		return
	}

	var userProjects []types.ProjectWithStats
	for _, p := range projectsWithStats {
        stats, err := server.store.GetProjectStats(ctx, util.ToNullInt4(int32(p.ID)))
        if err != nil {
            log.Printf("Error getting stats for project %d: %v", p.ID, err)
            continue
        }

		project := db.Project {
			ID: p.ID,
			Title: p.Title,
			Description: p.Description,
			StartDate: p.StartDate,
			EndDate: p.EndDate,
			CreatedBy: p.CreatedBy,
		}

		projectWithStats := types.ProjectWithStats{
			Project: project,
			TaskCount: stats.TaskCount,
			TotalTimeSpent: util.ToNullInt8(stats.TotalTimeSpent),
		}

		if p.CreatedBy.Int32 == userID {
			userProjects = append(userProjects, projectWithStats)
			continue
		}

		participants, err := server.store.ListProjectParticipants(ctx, p.ID)
		if err != nil {
			continue
		}

		for _, participant := range participants {
			if participant.ID == userID {
				userProjects = append(userProjects, projectWithStats)
				break
			}
		}
	}

	ctx.HTML(http.StatusOK, "projects.html", gin.H{
		"projects": userProjects,
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

func (server *Server) addProjectParticipant(ctx *gin.Context) {
	projectID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		UserID int64  `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.AddProjectParticipantParams{
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}

	participant, err := server.store.AddProjectParticipant(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, participant)
}

func (server *Server) removeProjectParticipant(ctx *gin.Context) {
	projectID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	userID, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	arg := db.RemoveProjectParticipantParams{
		ProjectID: projectID,
		UserID:    userID,
	}

	err = server.store.RemoveProjectParticipant(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (server *Server) showCreateProjectForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "create_project_modal.html", nil)
}

func (server *Server) createProject(ctx *gin.Context) {
	userID := getUserIDFromToken(ctx)

	arg := db.CreateProjectParams{
		Title:       ctx.PostForm("title"),
		Description: ctx.PostForm("description"),
		StartDate:   util.ParseDate(ctx.PostForm("start_date")),
		EndDate:     util.ParseDate(ctx.PostForm("end_date")),
		CreatedBy:   util.ToNullInt4(userID),
	}

	project, err := server.store.CreateProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "project_card.html", gin.H{
		"project": project,
	})
}
