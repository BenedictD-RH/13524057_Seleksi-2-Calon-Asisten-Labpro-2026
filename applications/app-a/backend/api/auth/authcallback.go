package auth

import (
	"app-a-backend/api/models"
	"app-a-backend/api/utility"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var local_session_expr = 24 * time.Hour

type AccessTokenPayload struct {
	Token string `json:"access_token" binding:"required"`
}

type UserInfoPayload struct {
	CentralSessionID datatypes.BinUUID   `json:"central_session_id" binding:"required"`
	UserID           datatypes.BinUUID   `json:"user_id" binding:"required"`
	Name             string              `json:"name" binding:"required"`
	Email            string              `json:"email" binding:"required"`
	Groups           []datatypes.BinUUID `json:"groups" binding:"required"`
}

func CreateLocalSession(user_info UserInfoPayload) (*models.LocalSession, string) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, ""
	}
	session_token, err := utility.CryptoRandString(32)
	if err != nil {
		return nil, ""
	}
	session_token_hash := utility.HashToken(session_token)

	return &models.LocalSession{
		ID:               datatypes.BinUUID(newUUID),
		SessionTokenHash: session_token_hash,
		ExternalUserId:   user_info.UserID,
		CentralSessionId: user_info.CentralSessionID,
		Status:           "Active",
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(local_session_expr),
	}, session_token
}

func CreateNewProfileCache(user_info UserInfoPayload) *models.ProfileCache {
	jsonData, err := json.Marshal(user_info.Groups)
	if err != nil {
		return nil
	}
	return &models.ProfileCache{
		ExternalUserId: user_info.UserID,
		Name:           user_info.Name,
		Email:          user_info.Email,
		GroupsList:     jsonData,
		SyncedAt:       time.Now(),
		CreatedAt:      time.Now(),
	}
}

func CreateUpdatedProfileCache(user_info UserInfoPayload) *models.ProfileCache {
	jsonData, err := json.Marshal(user_info.Groups)
	if err != nil {
		return nil
	}
	return &models.ProfileCache{
		Name:       user_info.Name,
		Email:      user_info.Email,
		GroupsList: jsonData,
		SyncedAt:   time.Now(),
	}
}

func UnauthorizedResponse(c *gin.Context, msg string) {
	newUUID, _:= uuid.NewRandom()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":     "INVALID_GRANT",
			"message":   msg,
			"requestId": newUUID.String(),
		},
	})
}

func InternalServerErrorResponse(c *gin.Context) {
	newUUID, _:= uuid.NewRandom()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":      "SERVER_ERROR",
			"message":   "Authorization request failed due to a server error",
			"requestId": newUUID.String(),
		},
	})
}

func AuthCallbackRequest(c *gin.Context, db *gorm.DB) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		UnauthorizedResponse(c, "Code or State not found! Please retry login!")
		return
	}
	auth_server_url := os.Getenv("AUTH_SERVER_URL")
	client_secret := os.Getenv("CLIENT_SECRET")
	client_id := os.Getenv("CLIENT_ID")
	if client_secret == "" || auth_server_url == "" || client_id == "" {
		InternalServerErrorResponse(c)
		return
	}

	var codeVerifierModels []models.CodeVerifier

	db.Where("state = ?", state).Find(&codeVerifierModels)

	if len(codeVerifierModels) != 1 {
		UnauthorizedResponse(c, "Code verifier not found! Please retry login!")
		return
	}
	
	payload := map[string]string{"code": code,
		"code_verifier": codeVerifierModels[0].CodeVerifier,
		"client_secret": client_secret,
		"redirect_uri": "https://" + c.Request.Host + c.FullPath()}
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/token", auth_server_url), bytes.NewReader(jsonData))

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		InternalServerErrorResponse(c)
		return
	}

	var tokenPayload AccessTokenPayload

	if err = json.NewDecoder(resp.Body).Decode(&tokenPayload); err != nil {
		InternalServerErrorResponse(c)
		return
	}
	resp.Body.Close()

	payload = map[string]string{"access_token": tokenPayload.Token}
	jsonData, _ = json.Marshal(payload)
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/userinfo", auth_server_url), bytes.NewReader(jsonData))

	req.Header.Set("Content-Type", "application/json")
	client = &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)

	if err != nil {
		InternalServerErrorResponse(c)
		return
	}

	var user_info UserInfoPayload
	if err = json.NewDecoder(resp.Body).Decode(&user_info); err != nil {
		InternalServerErrorResponse(c)
		return
	}
	resp.Body.Close()

	var profile_cache models.ProfileCache
	var profiles []models.ProfileCache
	db.Where("external_user_id = ?", user_info.UserID).Find(&profiles)
	if len(profiles) > 0 {
		updated_profile := CreateUpdatedProfileCache(user_info)
		if updated_profile == nil {
			InternalServerErrorResponse(c)
			return
		}
		db.Model(&models.ProfileCache{}).Where("external_user_id = ?", user_info.UserID).Updates(updated_profile)
	} else {
		if profile_pointer := CreateNewProfileCache(user_info); profile_pointer != nil {
			profile_cache = *profile_pointer
		} else {
			InternalServerErrorResponse(c)
			return
		}

		if err = db.Create(&profile_cache).Error; err != nil {
			InternalServerErrorResponse(c)
			return
		}
	}

	local_session, session_token := CreateLocalSession(user_info)
	if local_session == nil {
		InternalServerErrorResponse(c)
		return
	}
	if err = db.Create(&local_session).Error; err != nil {
		InternalServerErrorResponse(c)
		return
	}

	c.SetCookie("local_ssid", session_token, int(local_session_expr.Seconds()), "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{
		"redirect" : os.Getenv("FRONTEND_URI"),
	})
}
