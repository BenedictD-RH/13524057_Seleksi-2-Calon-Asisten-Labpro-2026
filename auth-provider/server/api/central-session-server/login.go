package centralsessionserver

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,gte=8"`
}

const session_exp_duration time.Duration = 24 * time.Hour

func UnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":      "INVALID_GRANT",
			"message":   "Authorization request is invalid",
			"requestId": requestid.Get(c),
		},
	})
}

func InternalServerResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":      "SERVER_ERROR",
			"message":   "Authorization request failed due to a server error",
			"requestId": requestid.Get(c),
		},
	})
}

func GenerateSession(user models.User) (models.SSOSession, error) {
	session_token, err := utility.CryptoRandString(16);
	if (err != nil) {
		return models.SSOSession{}, err
	}
	token_hash, err := utility.HashString(session_token)
	if (err != nil) {
		return models.SSOSession{}, err
	}
	return models.SSOSession{
		UserId: user.Id,
		SessionTokenHash: token_hash,
		Status: "Active",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(session_exp_duration),
	}, nil
}

func LoginRequest(c *gin.Context, db *gorm.DB) {
	var loginPayload LoginPayload

	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		UnauthorizedResponse(c)
		return
	}
	var users []models.User
	db.Where("email = ?", loginPayload.Email).Find(&users)

	if len(users) <= 0 {
		UnauthorizedResponse(c)
		return
	}

	if !users[0].IsActive() {
		UnauthorizedResponse(c)
		return
	}

	if !utility.CheckHash(loginPayload.Password, users[0].PasswordHash) {
		UnauthorizedResponse(c)
		return
	}
	session, err := GenerateSession(users[0])
	if (err != nil) {
		InternalServerResponse(c)
		return;
	}

	if err := db.Create(&session).Error; err != nil {
		InternalServerResponse(c)
		return;
	}

	c.SetCookie("ssid", session.Id.String(), int(session_exp_duration.Seconds()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login Successful",
	})
}
