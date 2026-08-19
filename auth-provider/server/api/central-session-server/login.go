package centralsessionserver

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,gte=8"`
}

const session_exp_duration time.Duration = 24 * time.Hour

func UnauthorizedLoginResponse(c *gin.Context, message string) {
	newUUID, _ := uuid.NewRandom()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":      "LOGIN_FAILED",
			"message":   message,
			"requestId": newUUID.String(),
		},
	})
}

func InternalServerResponse(c *gin.Context) {
	newUUID, _ := uuid.NewRandom()
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":      "SERVER_ERROR",
			"message":   "Authorization request failed due to a server error",
			"requestId": newUUID.String(),
		},
	})
}

func GenerateSession(user models.User) (models.SSOSession, string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.SSOSession{}, "", err
	}
	session_token, err := utility.CryptoRandString(16)
	if err != nil {
		return models.SSOSession{}, "", err
	}
	token_hash := utility.HashToken(session_token)
	return models.SSOSession{
		ID:               datatypes.BinUUID(newUUID),
		UserId:           user.ID,
		SessionTokenHash: token_hash,
		Status:           "Active",
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(session_exp_duration),
	}, session_token, nil
}

func LoginRequest(c *gin.Context, db *gorm.DB) {
	var loginPayload LoginPayload

	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		UnauthorizedLoginResponse(c, "Invalid login form!")
		return
	}
	var users []models.User
	db.Where("email = ?", loginPayload.Email).Find(&users)

	if len(users) <= 0 {
		UnauthorizedLoginResponse(c, "Wrong username or password!")
		return
	}

	if !users[0].IsActive() {
		UnauthorizedLoginResponse(c, "User account is inactive!")
		return
	}

	if !utility.CheckPasswordHash(loginPayload.Password, users[0].PasswordHash) {
		UnauthorizedLoginResponse(c, "Wrong username or password!")
		return
	}

	session, session_token, err := GenerateSession(users[0])
	if err != nil {
		InternalServerResponse(c)
		return
	}
	result := db.Create(&session)
	if err := result.Error; err != nil {
		InternalServerResponse(c)
		return
	}

	db.Model(&models.SSOSession{}).
	   Where("user_id = ? AND id != ? AND status = ?", users[0].ID, session.ID, "Active").
	   Updates(models.SSOSession{Status: "Revoked", RevokedReason: "new_session_created"})

	c.SetCookie("ssid", session_token, int(session_exp_duration.Seconds()), "/", "", false, true)
	AuthResponse(c, db, &session)
}
