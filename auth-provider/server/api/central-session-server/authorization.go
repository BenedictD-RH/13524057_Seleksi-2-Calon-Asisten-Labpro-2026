package centralsessionserver

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var auth_code_expr_duration = 2 * time.Minute

func GenerateAuthCode(code_challenge, redirect_uri string, session models.SSOSession, app models.Application, c *gin.Context, db *gorm.DB) (*models.AuthorizationCode, string) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		InternalServerResponse(c)
		return nil, ""
	}
	
	auth_code, err := utility.CryptoRandString(64)
	if err != nil {
		InternalServerResponse(c)
		return nil, ""
	}
	code_hash := utility.HashToken(auth_code)

	return &models.AuthorizationCode{
		ID:            datatypes.BinUUID(newUUID),
		CodeHash:      code_hash,
		CodeChallenge: code_challenge,
		UserId:        session.UserId,
		ApplicationId: app.ID,
		SsoSessionId:  session.ID,
		RedirectUri:   redirect_uri,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(auth_code_expr_duration),
	}, auth_code
}

func UnauthorizedResponse(c *gin.Context, msg string) {
	newUUID, _ := uuid.NewRandom()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":      "INVALID_GRANT",
			"message":   msg,
			"requestId": newUUID.String(),
		},
	})
}

func AuthorizeRequest(c *gin.Context, db *gorm.DB) {
	session_token, err := c.Cookie("ssid")
	if err != nil {
		RedirectToLoginPage(c)
		fmt.Println(err)
		return
	}
	var sessions []models.SSOSession
	db.Where("status = ? AND session_token_hash = ?", "Active", utility.HashToken(session_token)).Find(&sessions)
	if len(sessions) != 1 {
		RedirectToLoginPage(c)
		return
	}
	sessions[0].UpdateStatus(db)
	if (!sessions[0].IsValid()) {
		RedirectToLoginPage(c)
		return
	}
	session := &sessions[0]
	session.MarkActivity(db)

	client_id := c.Query("client_id")
	redirect_uri := c.Query("redirect_uri")
	code_challenge := c.Query("code_challenge")
	state := c.Query("state")
	if state == "" || client_id == "" || redirect_uri == "" || code_challenge == "" {
		UnauthorizedResponse(c, "Missing query parameters")
		return
	}
	auth_portal_uri := os.Getenv("AUTH_PORTAL_URI")

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/session?client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s",
		auth_portal_uri, client_id, redirect_uri, state, code_challenge))
}

func PickSessionRequest(c *gin.Context, db *gorm.DB) {
	session_token, err := c.Cookie("ssid")
	if err != nil {
		RedirectToLoginPage(c)
		fmt.Println(err)
		return
	}
	var sessions []models.SSOSession
	db.Where("status = ? AND session_token_hash = ?", "Active", utility.HashToken(session_token)).Find(&sessions)
	if len(sessions) != 1 {
		RedirectToLoginPage(c)
		return
	}
	sessions[0].UpdateStatus(db)
	if (!sessions[0].IsValid()) {
		RedirectToLoginPage(c)
		return
	}
	session := &sessions[0]
	AuthResponse(c, db, session)
}

func AuthResponse(c *gin.Context, db *gorm.DB, session *models.SSOSession) {
	client_id := c.Query("client_id")
	redirect_uri := c.Query("redirect_uri")
	code_challenge := c.Query("code_challenge")
	state := c.Query("state")
	if state == "" || client_id == "" || redirect_uri == "" || code_challenge == "" {
		UnauthorizedResponse(c, "Missing query parameters")
		return
	}

	var apps []models.Application
	db.Where("client_id = ?", client_id).Find(&apps)
	if len(apps) <= 0 {
		UnauthorizedResponse(c, "Invalid Client ID")
		return
	}

	if (apps[0].Status != "Active") {
		UnauthorizedResponse(c, "Application is not active")
		return
	}

	var redirect_uris []models.ApplicationRedirectURI
	db.Where("redirect_uri = ? AND application_id = ?", redirect_uri, apps[0].ID).Find(&redirect_uris)
	if len(redirect_uris) <= 0 {
		UnauthorizedResponse(c, "Redirect URI not registered")
		return
	}

	user := models.User{ID: session.UserId}
	if (!user.IsAuthorized(apps[0].ID, db)) {
		UnauthorizedResponse(c, "Policy Denied")
		return
	}

	if (!session.IsValid()) {
		UnauthorizedResponse(c, "Invalid session")
		return
	}

	auth_code_model, auth_code := GenerateAuthCode(code_challenge, redirect_uri, *session, apps[0], c, db)
	if auth_code_model == nil && auth_code == "" {
		UnauthorizedResponse(c, "Invalid Authorization Code")
		return
	}

	if err := db.Create(auth_code_model).Error; err != nil {
		InternalServerResponse(c)
		return
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s?code=%s&state=%s", redirect_uri, auth_code, state))
}

func RedirectToLoginPage(c *gin.Context) {
	client_id := c.Query("client_id")
	redirect_uri := c.Query("redirect_uri")
	code_challenge := c.Query("code_challenge")
	state := c.Query("state")
	if state == "" || client_id == "" || redirect_uri == "" || code_challenge == "" {
		UnauthorizedResponse(c, "Missing query parameters")
		return
	}
	auth_portal_uri := os.Getenv("AUTH_PORTAL_URI")

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/login?client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s",
		auth_portal_uri, client_id, redirect_uri, state, code_challenge))
}

func ReturnToClientPage(c *gin.Context, db *gorm.DB) {
	client_id := c.Query("client_id")
	if client_id == ""  {
		UnauthorizedResponse(c, "No Client Queried")
		return
	}

	var app models.Application
	db.Model(&models.Application{}).Where("client_id = ?", client_id).First(&app)

	c.JSON(http.StatusOK, gin.H{
		"client_url" : app.LaunchUrl,
	})
}

func AuthorizeAdminRequest(c *gin.Context, db *gorm.DB) {
	session_token, err := c.Cookie("ssid")
	if err != nil {
		UnauthorizedResponse(c, "Invalid session")
		return
	}
	var sessions []models.SSOSession
	db.Where("status = ? AND session_token_hash = ?", "Active", utility.HashToken(session_token)).Find(&sessions)
	if len(sessions) != 1 {
		UnauthorizedResponse(c, "Invalid session")
		return
	}
	sessions[0].UpdateStatus(db)
	if (!sessions[0].IsValid()) {
		UnauthorizedResponse(c, "Invalid session")
		return
	}
	session := &sessions[0]

	var user models.User
	db.Model(&models.User{}).Where("id = ?", session.UserId).First(&user)

	var adminGroup models.Group
	db.Model(&models.Group{}).Where("name = ?", "administrators").First(&adminGroup)

	var userGroup []models.UserGroup
	db.Model(&models.UserGroup{}).Where("user_id = ? AND group_id = ?", user.ID, adminGroup.ID).First(&userGroup)

	if len(userGroup) == 0 {
		UnauthorizedResponse(c, "User is unauthorized")
		return
	}

	if !user.IsActive() {
		UnauthorizedResponse(c, "User is inactive")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : "authorized",
	})
}