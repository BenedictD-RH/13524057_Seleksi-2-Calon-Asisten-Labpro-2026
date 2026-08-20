package centralsessionserver

import (
	eventpublisher "auth-provider-server/api/event-publisher"
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func LogoutRequest(c *gin.Context, db, mq *gorm.DB) {
	session_token, err := c.Cookie("ssid")
	if err != nil {
		UnauthorizedResponse(c, "No session found")
		return
	}

	var sessions []models.SSOSession

	db.Model(&models.SSOSession{}).
	   Where("session_token_hash = ? AND status = ?", utility.HashToken(session_token), "Active").
	   Find(&sessions)
	
	if len(sessions) != 1 {
		models.AuditEvent("logout_attempt", "failed", nil, nil, nil, c, db)
		InternalServerResponse(c)
		return
	}

	eventpublisher.PublishEvent(sessions[0].UserId, &sessions[0].ID, nil, "SessionRevoked", "user_logout", mq)

	db.Model(&models.SSOSession{}).
	   Where("session_token_hash = ?", utility.HashToken(session_token)).
	   Updates(models.SSOSession{Status: "Revoked", RevokedReason: "user_logout"})

	c.SetCookie("ssid", "", -1, "/", "", false, true)
	models.AuditEvent("logout_attempt", "success", &sessions[0].UserId, nil, &sessions[0].ID, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout successful",
	})
}