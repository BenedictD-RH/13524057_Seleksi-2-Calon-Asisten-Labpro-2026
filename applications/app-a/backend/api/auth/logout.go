package auth

import (
	"app-a-backend/api/models"
	"app-a-backend/api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func LogoutRequest(c *gin.Context, db *gorm.DB) {
	local_session_token, err := c.Cookie("local_ssid")
	if err != nil {
		UnauthorizedResponse(c)
		return
	}

	db.Model(&models.LocalSession{}).
	   Where("session_token_hash = ?", utility.HashToken(local_session_token)).
	   Updates(models.LocalSession{Status: "Revoked", RevokedReason: "user_logout"})

	c.SetCookie("local_ssid", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout successful",
	})
}