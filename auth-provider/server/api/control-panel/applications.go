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
	Name           string `json:"name" binding:"required"`
	ClientID       string `json:"client_id" binding:"required"`
	ClientSecret   string `json:"client_secret" binding:"required,gte=8"`
	LaunchURL      string `json:"launch_url"`
	LogoutNotifURL string `json:"logout_notification_url"`
}

type UpdateAppPayload struct {
	Id       datatypes.BinUUID `json:"id" binding:"required"`
	Name           string `json:"name"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret" binding:"omitempty,gte=8"`
	LaunchURL      string `json:"launch_url"`
	LogoutNotifURL string `json:"logout_notification_url"`
}

func (payload UpdateAppPayload) GetDBStruct() (models.Application, models.Application, error) {
	var hashedPass string = ""
	var err error = nil
	if payload.ClientSecret != "" {
		hashedPass, err = utility.HashPassword(payload.ClientSecret)
	}

	return models.Application{
			ID: payload.Id,
		},
		models.Application{
			Name:         	  payload.Name,
			ClientId:         payload.ClientID,
			ClientSecretHash: hashedPass,
			LaunchUrl:        payload.LaunchURL,
			LogoutNotificationUrl: payload.LogoutNotifURL,
		}, err
}

func (payload PKeyPayload) GetAppModel() models.Application {
	return models.Application{ID: payload.Id}
}

type RegisterAppURIPayload struct {
	ClientID    string `json:"client_id" binding:"required"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
}



func (payload RegisterAppPayload) GetDBStruct() (models.Application, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.Application{}, err
	}
	client_secret_hash, err := utility.HashPassword(payload.ClientSecret)
	if err != nil {
		return models.Application{}, err
	}

	return models.Application{
		ID:                    datatypes.BinUUID(newUUID),
		Name:                  payload.Name,
		ClientId:              payload.ClientID,
		ClientSecretHash:      client_secret_hash,
		Status:                "Active",
		LaunchUrl:             payload.LaunchURL,
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

func RegisterAppURIRequest(c *gin.Context, db *gorm.DB) {
	var uriPayload RegisterAppURIPayload

	if err := c.ShouldBindJSON(&uriPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to register app: " + err.Error(),
		})
		return
	}
	
	var apps []models.Application
	db.Find("client_id = ?", uriPayload.ClientID).Find(&apps)
	if len(apps) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to register app: client_id not registered",
		})
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to register app: " + err.Error(),
		})
		return
	}

	appUriModel := models.ApplicationRedirectURI{
		ID:                    datatypes.BinUUID(newUUID),
		ApplicationId: apps[0].ID,
		RedirectUri: uriPayload.RedirectURI,
	}

	if err := db.Create(&appUriModel).Error; err != nil {
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


func GetAllApplicationsRequest(c *gin.Context, db *gorm.DB) {
	var apps []models.Application

	if err := db.Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get app data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, apps)
}

func UpdateAppRequest(c *gin.Context, db *gorm.DB) {
	var appPayload UpdateAppPayload

	if err := c.ShouldBindJSON(&appPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update app: " + err.Error(),
		})
		return
	}

	appQuery, updatedAppModel, err := appPayload.GetDBStruct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update app: " + err.Error(),
		})
		return
	}

	if err := db.Model(&appQuery).Updates(updatedAppModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update app: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "App successfully updated",
	})
}

func RemoveAppRequest(c *gin.Context, db *gorm.DB) {
	var appPKey PKeyPayload

	if err := c.ShouldBindJSON(&appPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove app: " + err.Error(),
		})
		return
	}

	appModel := appPKey.GetAppModel()
	if err := db.Delete(&appModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove app: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "App successfully removed",
	})
}