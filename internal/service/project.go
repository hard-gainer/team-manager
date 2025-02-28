package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

		project := db.Project{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			StartDate:   p.StartDate,
			EndDate:     p.EndDate,
			CreatedBy:   p.CreatedBy,
		}

		projectWithStats := types.ProjectWithStats{
			Project:        &project,
			TaskCount:      stats.TaskCount,
			TotalTimeSpent: util.ToNullInt8(stats.TotalTimeSpent),
		}

		// if p.CreatedBy.Int32 == userID {
		// 	userProjects = append(userProjects, projectWithStats)
		// 	continue
		// }

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

	log.Printf("Found %d projects for user %d", len(userProjects), userID)

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

	user, err := server.store.GetEmployee(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg1 := db.AddProjectParticipantParams{
		ProjectID: project.ID,
		UserID:    int64(user.ID),
		Role:      user.Role,
	}

	_, err = server.store.AddProjectParticipant(ctx, arg1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "project_card.html", gin.H{
		"project": project,
	})
}

func (server *Server) showInviteMemberForm(ctx *gin.Context) {
	projectID := ctx.Param("id")
	ctx.HTML(http.StatusOK, "invite_member_modal.html", gin.H{
		"projectID": projectID,
	})
}

func (server *Server) inviteMember(ctx *gin.Context) {
	projectID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	email := ctx.PostForm("email")
	role := ctx.PostForm("role")

	if email == "" || role == "" {
		ctx.HTML(http.StatusBadRequest, "invite_result.html", gin.H{
			"error": "Email and role are required",
		})
		return
	}

	log.Printf("generating invitation link!")
	token := util.GenerateSecureToken(32)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	log.Printf("Creating invitation: projectID=%d, email=%s, role=%s, expiresAt=%v",
		projectID, email, role, expiresAt)

	arg := db.CreateProjectInvitationParams{
		ProjectID: projectID,
		Email:     email,
		Token:     token,
		Role:      role,
		ExpiresAt: util.ToTimestamp(expiresAt),
	}

	_, err = server.store.CreateProjectInvitation(ctx, arg)
	if err != nil {
		ctx.HTML(http.StatusOK, "invite_result.html", gin.H{
			"error": "Failed to create invitation",
		})
		return
	}

	inviteLink := fmt.Sprintf("%s/projects/join/%s", server.config.BaseURL, token)
	err = server.mailer.SendInvitation(email, inviteLink)
	if err != nil {
		ctx.HTML(http.StatusOK, "invite_result.html", gin.H{
			"error": "Failed to send invitation email",
		})
		return
	}

	log.Println("Success!!!")

	ctx.HTML(http.StatusOK, "invite_result.html", gin.H{
		"success": "Invitation sent successfully",
	})
}

func (server *Server) handleProjectInvitation(ctx *gin.Context) {
	log.Println("service.handleProjectInvitation")

	token := ctx.Param("token")
	if token == "" {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid invitation link",
		})
		return
	}

	userID := getUserIDFromToken(ctx)
	if userID == 0 {
		returnURL := fmt.Sprintf("/projects/join/%s", token)
		ctx.Redirect(http.StatusSeeOther, "/login?return_url="+returnURL)
		return
	}

	invitation, err := server.store.GetProjectInvitation(ctx, token)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid or expired invitation",
		})
		return
	}

	if time.Now().After(invitation.ExpiresAt.Time) {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invitation has expired",
		})
		return
	}

	if invitation.AcceptedAt.Valid {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invitation has already been used",
		})
		return
	}

	user, err := server.store.GetEmployee(ctx, userID)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to verify user",
		})
		return
	}

	log.Printf("user.Email: %s, invitation.Email: %s\n", user.Email, invitation.Email)

	if user.Email != invitation.Email {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "This invitation was sent to a different email address",
		})
		return
	}

	arg := db.AddProjectParticipantParams{
		ProjectID: invitation.ProjectID,
		UserID:    int64(userID),
		Role:      invitation.Role,
	}

	_, err = server.store.AddProjectParticipant(ctx, arg)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to add you to the project",
		})
		return
	}

	_, err = server.store.AcceptProjectInvitation(ctx, token)
	if err != nil {
		log.Printf("Failed to mark invitation as accepted: %v", err)
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/dashboard/%d", invitation.ProjectID))
}
