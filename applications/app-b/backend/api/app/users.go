package app

import (
	"app-a-backend/api/models"
	"app-a-backend/api/utility"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetUserDataRequest(c *gin.Context, db *gorm.DB) {
	local_session_token, err := c.Cookie(os.Getenv("CLIENT_ID") + "_local_ssid")
	if err != nil {
		NoContentResponse(c)
		return
	}
	var sessions []models.LocalSession
	db.Where("status = ? AND session_token_hash = ?", "Active", utility.HashToken(local_session_token)).Find(&sessions)
	if len(sessions) != 1 {
		NoContentResponse(c)
		return
	}
	local_session := sessions[0]
	local_session.UpdateStatus(db)
	if (!local_session.IsValid()) {
		NoContentResponse(c)
		return
	}

	local_session.MarkActivity(db)
	var profile []models.ProfileCache
	db.Where("external_user_id = ?", local_session.ExternalUserId).Find(&profile)

	if len(profile) != 1 {
		NoContentResponse(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  profile[0].Name,
		"email": profile[0].Email,
	})
}
