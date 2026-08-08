package controlpanel

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RegisterAppPayload struct {
	Name     string `json:"name" binding:"required"`
	ClientID string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectURL string `json:"redirect_url"`
	LogoutNotifURL string `json:"logout_notification_url"`
}

func (payload RegisterAppPayload) GetDBStruct() (models.Application, error) {
	newUUID, err := uuid.NewRandom()
    if err != nil {
        return models.Application{}, err
    }
	client_secret_hash, err := utility.HashString(payload.ClientSecret)
	if (err != nil) {
		return models.Application{}, err
	}
	
	return models.Application{
		ID: datatypes.BinUUID(newUUID),
		Name: payload.Name,
		ClientId: payload.ClientID,
		ClientSecretHash: client_secret_hash,
		Status: "Active",
		RedirectUrl: payload.RedirectURL,
		LogoutNotificationUrl: payload.LogoutNotifURL,
	}, nil
}

func RegisterApplicationRequest(c *gin.Context, db *gorm.DB) {
	var appPayload RegisterAppPayload 

	if err := c.ShouldBindJSON(&appPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to register app: " + err.Error(),
		})
		return
	}
	appModel, err := appPayload.GetDBStruct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to register app: " + err.Error(),
		})
		return
	}

	if err := db.Create(&appModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to register app: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "App successfully registered",
	})
}