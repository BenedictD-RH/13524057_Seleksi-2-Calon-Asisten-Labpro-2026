package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserGroupPayload struct {
	UserId datatypes.BinUUID	`json:"user_id" binding:"required"`
	GroupId datatypes.BinUUID	`json:"group_id" binding:"required"`
}

type UserGroup struct {
	Id datatypes.BinUUID 		`json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	UserId datatypes.BinUUID	`json:"user_id"`
	GroupId datatypes.BinUUID	`json:"group_id"`
	CreatedAt time.Time         `json:"created_at"`
}

func (payload UserGroupPayload) GetDBStruct() (UserGroup) {
	return UserGroup{
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

	userGroupModel := userGroupPayload.GetDBStruct()

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

	userGroupModel := userGroupPayload.GetDBStruct()

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