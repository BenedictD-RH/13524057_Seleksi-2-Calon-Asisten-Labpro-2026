package centralsessionserver

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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

func GenerateSession(user models.User) (models.SSOSession, string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.SSOSession{}, "", err
	}
	session_token, err := utility.CryptoRandString(16)
	if err != nil {
		return models.SSOSession{}, "", err
	}
	token_hash, err := utility.HashPassword(session_token)
	if err != nil {
		return models.SSOSession{}, "", err
	}
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
		UnauthorizedResponse(c)
		fmt.Println("Payload missing")
		return
	}
	var users []models.User
	db.Where("email = ?", loginPayload.Email).Find(&users)

	if len(users) <= 0 {
		UnauthorizedResponse(c)
		fmt.Println("User not found")
		return
	}

	if !users[0].IsActive() {
		UnauthorizedResponse(c)
		fmt.Println("User not active")
		return
	}

	if !utility.CheckPasswordHash(loginPayload.Password, users[0].PasswordHash) {
		UnauthorizedResponse(c)
		fmt.Println("Wrong password")
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

	c.SetCookie("ssid", session_token, int(session_exp_duration.Seconds()), "/", "", false, true)
	AuthResponse(c, db, &session)

}
