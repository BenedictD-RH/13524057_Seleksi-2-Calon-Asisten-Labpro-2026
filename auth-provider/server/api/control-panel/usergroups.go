package controlpanel

import (
	"net/http"
	"github.com/google/uuid"
	"auth-provider-server/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserGroupPayload struct {
	UserId datatypes.BinUUID	`json:"user_id" binding:"required"`
	GroupId datatypes.BinUUID	`json:"group_id" binding:"required"`
}



func (payload UserGroupPayload) GetDBStructCreate() (models.UserGroup, error) {
	newUUID, err := uuid.NewRandom()
    if err != nil {
        return models.UserGroup{}, err
    }
	return models.UserGroup{
		ID:	datatypes.BinUUID(newUUID),
		UserId: payload.UserId,
		GroupId: payload.GroupId,
	}, nil
}

func (payload UserGroupPayload) GetDBStructDelete() (models.UserGroup) {
	return models.UserGroup{
		UserId: payload.UserId,
		GroupId: payload.GroupId,
	}
}

func AddUserToGroup(c *gin.Context, db *gorm.DB) {
	var userGroupPayload UserGroupPayload

	if err := c.ShouldBindJSON(&userGroupPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to add user to group: " + err.Error(),
		})
		return
	}

	userGroupModel, err := userGroupPayload.GetDBStructCreate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to add user to group: " + err.Error(),
		})
		return
	}

	if err := db.Create(&userGroupModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to add user to group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Successfully added user to group",
	})
}

func RemoveUserFromGroup(c *gin.Context, db *gorm.DB) {
	var userGroupPayload UserGroupPayload

	if err := c.ShouldBindJSON(&userGroupPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	userGroupModel := userGroupPayload.GetDBStructDelete()

	if err := db.Delete(&userGroupModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully removed user to group",
	})
}