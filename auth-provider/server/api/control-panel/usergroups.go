package controlpanel

import (
	"auth-provider-server/api/models"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserGroupPayload struct {
	UserId  datatypes.BinUUID `json:"user_id" binding:"required"`
	GroupId datatypes.BinUUID `json:"group_id" binding:"required"`
}

type UserFromGroupPayload struct {
	UserId   datatypes.BinUUID `json:"user_id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Status   string			   `json:"status"`
	PlacedAt time.Time         `json:"placed_at"`
}

type GroupFromUserPayload struct {
	GroupId     datatypes.BinUUID `json:"group_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	PlacedAt    time.Time         `json:"placed_at"`
}

func (payload UserGroupPayload) GetDBStructCreate() (models.UserGroup, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.UserGroup{}, err
	}
	return models.UserGroup{
		ID:      datatypes.BinUUID(newUUID),
		UserId:  payload.UserId,
		GroupId: payload.GroupId,
	}, nil
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

	db.Where("user_id = ? AND group_id = ?", userGroupPayload.UserId, userGroupPayload.GroupId).Delete(&models.UserGroup{})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully removed user to group",
	})
}

func GetUsersFromGroupRequest(c *gin.Context, db *gorm.DB) {
	var usersFromGroup []UserFromGroupPayload

	var groupPKey PKeyPayload

	if err := c.ShouldBindJSON(&groupPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	db.Model(&models.UserGroup{}).
		Select("user_groups.user_id, users.name, users.email, users.status, user_groups.created_at AS placed_at").
		Joins("LEFT JOIN users ON users.id = user_groups.user_id").
		Where("group_id = ?", groupPKey.Id).
		Find(&usersFromGroup)

	c.JSON(http.StatusOK, usersFromGroup)
}

func GetGroupsFromUserRequest(c *gin.Context, db *gorm.DB) {
	var groupFromUsers []GroupFromUserPayload

	var userPKey PKeyPayload

	if err := c.ShouldBindJSON(&userPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	db.Model(&models.UserGroup{}).
		Select("user_groups.group_id, groups_.name, groups_.description, user_groups.created_at AS placed_at").
		Joins("LEFT JOIN groups_ ON user_groups.group_id = groups_.id").
		Where("user_groups.user_id = ?", userPKey.Id).
		Find(&groupFromUsers)

	c.JSON(http.StatusOK, groupFromUsers)
}

func GetUsersNotInGroupRequest(c *gin.Context, db *gorm.DB) {
	var users []UserCompact

	var groupPKey PKeyPayload

	if err := c.ShouldBindJSON(&groupPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	usersInGroup := db.Model(&models.User{}).
					   Select("users.id").
					   Joins("RIGHT JOIN user_groups ON user_groups.user_id = users.id").
					   Where("user_groups.group_id = ?", groupPKey.Id)
	
	db.Model(&models.User{}).Select("id, name, email").Where("id NOT IN (?)", usersInGroup).Find(&users)

	c.JSON(http.StatusOK, users)
}

func GetGroupsUserIsNotInRequest(c *gin.Context, db *gorm.DB) {
	var groups []GroupCompact

	var userPKey PKeyPayload

	if err := c.ShouldBindJSON(&userPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	groupsUserIsIn := db.Model(&models.Group{}).
					   Select("groups_.id").
					   Joins("RIGHT JOIN user_groups ON user_groups.group_id = groups_.id").
					   Where("user_groups.user_id = ?", userPKey.Id)
	
	db.Model(&models.Group{}).Select("id, name").Where("id NOT IN (?)", groupsUserIsIn).Find(&groups)

	c.JSON(http.StatusOK, groups)
}

func GetUsersFromGroupFieldsRequest(c *gin.Context) {
	userFromGroup := UserFromGroupPayload{}

	t:= reflect.TypeOf(userFromGroup)
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"fields" : fieldNames,
	})
}

func GetGroupsFromUserFieldsRequest(c *gin.Context) {
	groupFromUser := GroupFromUserPayload{}

	t:= reflect.TypeOf(groupFromUser)
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"fields" : fieldNames,
	})
}
