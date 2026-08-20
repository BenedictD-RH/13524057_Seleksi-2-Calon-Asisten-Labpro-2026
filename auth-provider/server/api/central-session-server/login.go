package centralsessionserver

import (
	eventpublisher "auth-provider-server/api/event-publisher"
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

func GenerateSession(user models.User, c *gin.Context) (models.SSOSession, string, error) {
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
		IpAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(session_exp_duration),
	}, session_token, nil
}

func LoginRequest(c *gin.Context, db, mq *gorm.DB) {
	var loginPayload LoginPayload

	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		models.AuditEvent("login_attempt", "failed", nil, nil, nil, c, db)
		UnauthorizedLoginResponse(c, "Invalid login form!")
		return
	}
	var users []models.User
	db.Where("email = ?", loginPayload.Email).Find(&users)

	if len(users) <= 0 {
		models.AuditEvent("login_attempt", "failed", nil, nil, nil, c, db)
		UnauthorizedLoginResponse(c, "Wrong username or password!")
		return
	}

	if !users[0].IsActive() {
		models.AuditEvent("login_attempt", "failed", &users[0].ID, nil, nil, c, db)
		UnauthorizedLoginResponse(c, "User account is inactive!")
		return
	}

	if !utility.CheckPasswordHash(loginPayload.Password, users[0].PasswordHash) {
		models.AuditEvent("login_attempt", "failed", &users[0].ID, nil, nil, c, db)
		UnauthorizedLoginResponse(c, "Wrong username or password!")
		return
	}

	session, session_token, err := GenerateSession(users[0], c)
	if err != nil {
		models.AuditEvent("login_attempt", "failed", &users[0].ID, nil, nil, c, db)
		InternalServerResponse(c)
		return
	}
	result := db.Create(&session)
	if err := result.Error; err != nil {
		models.AuditEvent("login_attempt", "failed", &users[0].ID, nil, nil, c, db)
		InternalServerResponse(c)
		return
	}

	db.Model(&models.SSOSession{}).
	   Where("user_id = ? AND id != ? AND status = ?", users[0].ID, session.ID, "Active").
	   Updates(models.SSOSession{Status: "Revoked", RevokedReason: "new_session_created"})


	old_session_token, err := c.Cookie("ssid")
	if err == nil {
		var sessions []models.SSOSession

		db.Model(&models.SSOSession{}).
			Where("session_token_hash = ? AND status = ?", utility.HashToken(old_session_token), "Active").
			Find(&sessions)
		
		if len(sessions) != 1 {
			models.AuditEvent("login_attempt", "failed", &users[0].ID, nil, &session.ID, c, db)
			InternalServerResponse(c)
			return
		}

		eventpublisher.PublishEvent(sessions[0].UserId, &sessions[0].ID, nil, "SessionRevoked", "user_logout", mq)

		db.Model(&models.SSOSession{}).
			Where("session_token_hash = ?", utility.HashToken(old_session_token)).
			Updates(models.SSOSession{Status: "Revoked", RevokedReason: "user_logout"})
	}

	
	

	c.SetCookie("ssid", session_token, int(session_exp_duration.Seconds()), "/", "localhost", false, true)

	client_id := c.Query("client_id")
	redirect_uri := c.Query("redirect_uri")
	code_challenge := c.Query("code_challenge")
	state := c.Query("state")
	if state == "" && client_id == "" && redirect_uri == "" && code_challenge == "" {
		models.AuditEvent("login_attempt", "success", &users[0].ID, nil, &session.ID, c, db)
		c.JSON(http.StatusOK, gin.H{
			"status" : "authorized",
		})
		return
	}
	var app models.Application

	db.Model(&models.Application{}).Where("client_id = ?", client_id).Find(&app)

	models.AuditEvent("login_attempt", "success", &users[0].ID, &app.ID, &session.ID, c, db)

	AuthResponse(c, db, &session)
}
