package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CreateGroupPayload struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateGroupPayload struct {
	Id          datatypes.BinUUID `json:"id" binding:"required"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
}

type Group struct {
	Id          datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (Group) TableName() string {
	return "groups_"
}

func (payload CreateGroupPayload) GetDBStruct() Group {
	return Group{
		Name:        payload.Name,
		Description: payload.Description,
	}
}

func (payload UpdateGroupPayload) GetDBStruct() (Group, Group) {
	return Group{
		Id: payload.Id,
	},
	Group{
		Name:        payload.Name,
		Description: payload.Description,
	}
}

func (payload PKeyPayload) GetGroupModel() Group {
	return Group{Id: payload.Id}
}

func CreateGroupRequest(c *gin.Context, db *gorm.DB) {
	var groupPayload CreateGroupPayload

	if err := c.ShouldBindJSON(&groupPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create group: " + err.Error(),
		})
		return
	}

	groupModel := groupPayload.GetDBStruct()

	if err := db.Create(&groupModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Group successfully created",
	})
}

func GetAllGroupsRequest(c *gin.Context, db *gorm.DB) {
	var groups []Group

	if err := db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get group data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func UpdateGroupRequest(c *gin.Context, db *gorm.DB) {
	var groupPayload UpdateGroupPayload

	if err := c.ShouldBindJSON(&groupPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update group: " + err.Error(),
		})
		return
	}

	userGroup, updatedGroupModel := groupPayload.GetDBStruct()

	if err := db.Model(&userGroup).Updates(updatedGroupModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Group successfully updated",
	})
}

func DeleteGroupRequest(c *gin.Context, db *gorm.DB) {
	var userPKey PKeyPayload

	if err := c.ShouldBindJSON(&userPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete group: " + err.Error(),
		})
		return
	}

	userModel := userPKey.GetGroupModel()
	if err := db.Delete(&userModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Group successfully deleted",
	})
}
