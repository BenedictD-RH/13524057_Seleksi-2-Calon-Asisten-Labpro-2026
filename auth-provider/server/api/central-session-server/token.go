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

var token_expr_duration = 20 * time.Second

type TokenRequestPayload struct {
	Code         string `json:"code" binding:"required"`
	CodeVerifier string `json:"code_verifier" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
}

func GenerateAccessToken(c *gin.Context, auth_code *models.AuthorizationCode) (*models.AccessToken, string) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, ""
	}
	access_token, err := utility.CryptoRandString(64)
	if err != nil {
		return nil, ""
	}
	access_token_hash, err := utility.HashPassword(access_token)
	if err != nil {
		return nil, ""
	}

	return &models.AccessToken{
		ID:            datatypes.BinUUID(newUUID),
		TokenHash:     access_token_hash,
		UserId:        (*auth_code).UserId,
		ApplicationId: (*auth_code).ApplicationId,
		SsoSessionId:  (*auth_code).SsoSessionId,
		Status:        "Active",
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(token_expr_duration),
	}, access_token
}

func TokenRequest(c *gin.Context, db *gorm.DB) {
	var tokenPayload TokenRequestPayload

	if err := c.ShouldBindJSON(&tokenPayload); err != nil {
		UnauthorizedResponse(c)
		return
	}

	var auth_codes []models.AuthorizationCode
	db.Where("code_hash = ?", utility.HashToken(tokenPayload.Code)).Find(&auth_codes)
	if len(auth_codes) != 1 {
		UnauthorizedResponse(c)
		return
	}
	auth_code := &auth_codes[0]
	if !utility.CheckPasswordHash(tokenPayload.CodeVerifier, auth_code.CodeChallenge) {
		UnauthorizedResponse(c)
		return
	}
	var app models.Application
	var apps []models.Application
	db.Where("id = ?", auth_code.ApplicationId).Find(&apps)
	if len(apps) <= 0 {
		UnauthorizedResponse(c)
		return
	}
	app = apps[0]

	if !utility.CheckPasswordHash(tokenPayload.ClientSecret, app.ClientSecretHash) {
		UnauthorizedResponse(c)
		return
	}

	accessTokenModel, access_token := GenerateAccessToken(c, auth_code)
	if accessTokenModel == nil {
		InternalServerResponse(c)
		return
	}
	if err := db.Create(&accessTokenModel).Error; err != nil {
		InternalServerResponse(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":       "success",
		"message":      "Access Token successfully generated",
		"access_token": access_token,
	})
}
