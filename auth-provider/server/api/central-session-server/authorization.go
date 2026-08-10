package centralsessionserver

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var auth_code_expr_duration = 2 * time.Minute

func GenerateAuthCode(code_challenge, client_id, redirect_uri string, session models.SSOSession, c *gin.Context, db *gorm.DB) (*models.AuthorizationCode, string) {
	newUUID, err := uuid.NewRandom()
    if err != nil {
		InternalServerResponse(c)
        return nil, ""
    }
	var apps []models.Application
	db.Where("client_id = ?", client_id).Find(&apps)
	if len(apps) <= 0 {
		UnauthorizedResponse(c)
		return nil, ""
	}
	auth_code, err := utility.CryptoRandString(64)
	if err != nil {
		InternalServerResponse(c)
		return nil, ""
	}
	code_hash, err := utility.HashString(auth_code)
	if err != nil {
		InternalServerResponse(c)
		return nil, ""
	}

	return &models.AuthorizationCode{
		ID:			   datatypes.BinUUID(newUUID),
		CodeHash:      code_hash,
		CodeChallenge: code_challenge,
		UserId:        session.UserId,
		ApplicationId: apps[0].ID,
		SsoSessionId:  session.ID,
		RedirectUri:   redirect_uri,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(auth_code_expr_duration),
	}, auth_code
}

func FindSession(sessions *[]models.SSOSession, session_token string) *models.SSOSession {
	for i := 0; i < len(*sessions); i++ {
		if utility.CheckHash(session_token, (*sessions)[i].SessionTokenHash) {
			return &(*sessions)[i]
		}
	}
	return nil
}

func AuthorizeRequest(c *gin.Context, db *gorm.DB) {
	session_token, err := c.Cookie("ssid")
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"status":  "login page",
			"message": "Please login",
		})
		return
	}
	var sessions []models.SSOSession
	db.Where("status = ?", "Active").Find(&sessions)
	session := FindSession(&sessions, session_token)
	if session == nil {
		c.JSON(http.StatusCreated, gin.H{
			"status":  "login page",
			"message": "Please login",
		})
		return
	}

	AuthResponse(c, db, session)
}

func AuthResponse(c *gin.Context, db *gorm.DB, session *models.SSOSession) {
	client_id := c.Query("client_id")
	redirect_uri := c.Query("redirect_uri")
	code_challenge := c.Query("code_challenge")
	state := c.Query("state")
	if state == "" || client_id == "" || redirect_uri == "" || code_challenge == "" {
		fmt.Println("Missing parameters")
		UnauthorizedResponse(c)
		return
	}
	auth_code_model, auth_code := GenerateAuthCode(code_challenge, client_id, redirect_uri, *session, c, db)
	if auth_code_model == nil && auth_code == "" {
		return
	}

	if err := db.Create(auth_code_model).Error; err != nil {
		InternalServerResponse(c)
		return
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/auth/callback?code=%s&state=%s", redirect_uri, auth_code, state))
}
