package controlpanel

import (
	eventpublisher "auth-provider-server/api/event-publisher"
	"auth-provider-server/api/models"
	"net/http"
	"reflect"

	"github.com/google/uuid"

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

type GroupCompact struct {
	Id          datatypes.BinUUID `json:"id"`
	Name        string            `json:"name"`
}

func (payload CreateGroupPayload) GetDBStruct() (models.Group, error) {
	newUUID, err := uuid.NewRandom()
    if err != nil {
        return models.Group{}, err
    }
	return models.Group{
		ID:			 datatypes.BinUUID(newUUID),
		Name:        payload.Name,
		Description: payload.Description,
	}, nil
}

func (payload UpdateGroupPayload) GetDBStruct() (models.Group, models.Group) {
	return models.Group{
			ID: payload.Id,
		},
		models.Group{
			Name:        payload.Name,
			Description: payload.Description,
		}
}

func (payload PKeyPayload) GetGroupModel() models.Group {
	return models.Group{ID: payload.Id}
}

func CreateGroupRequest(c *gin.Context, db *gorm.DB) {
	var groupPayload CreateGroupPayload

	if err := c.ShouldBindJSON(&groupPayload); err != nil {
		models.AuditEvent("create_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create group: " + err.Error(),
		})
		return
	}

	groupModel, err := groupPayload.GetDBStruct()
	if err != nil {
		models.AuditEvent("create_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create group: " + err.Error(),
		})
		return
	}

	if err := db.Create(&groupModel).Error; err != nil {
		models.AuditEvent("create_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create group: " + err.Error(),
		})
		return
	}
	models.AuditEvent("create_group", "success", nil, nil, nil, c, db)
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Group successfully created",
	})
}

func GetAllGroupsRequest(c *gin.Context, db *gorm.DB) {
	var groups []models.Group

	if err := db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get group data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func GetGroupFields(c *gin.Context) {
	user := models.Group{}

	t:= reflect.TypeOf(user)
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"fields" : fieldNames,
	})
}

func UpdateGroupRequest(c *gin.Context, db *gorm.DB) {
	var groupPayload UpdateGroupPayload

	if err := c.ShouldBindJSON(&groupPayload); err != nil {
		models.AuditEvent("update_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update group: " + err.Error(),
		})
		return
	}

	userGroup, updatedGroupModel := groupPayload.GetDBStruct()

	if err := db.Model(&userGroup).Updates(updatedGroupModel).Error; err != nil {
		models.AuditEvent("update_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update group: " + err.Error(),
		})
		return
	}
	models.AuditEvent("update_group", "success", nil, nil, nil, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Group successfully updated",
	})
}

func DeleteGroupRequest(c *gin.Context, db *gorm.DB, mq *gorm.DB) {
	var groupPKey PKeyPayload

	if err := c.ShouldBindJSON(&groupPKey); err != nil {
		models.AuditEvent("delete_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete group: " + err.Error(),
		})
		return
	}

	OnGroupDelete(groupPKey.Id, db, mq)

	groupModel := groupPKey.GetGroupModel()
	if err := db.Delete(&groupModel).Error; err != nil {
		models.AuditEvent("delete_group", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete group: " + err.Error(),
		})
		return
	}

	models.AuditEvent("delete_group", "success", nil, nil, nil, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Group successfully deleted",
	})
}

func OnGroupDelete(group_id datatypes.BinUUID, db, mq *gorm.DB) {
	var appPKeys []PKeyPayload

	user_ids := db.Model(&models.UserGroup{}).
					Select("user_id AS id").
					Where("group_id = ?", group_id)
		
	db.Model(&models.ApplicationGroupPolicy{}).
	   Select("application_id AS id").
	   Where("group_id = ? AND effect = ?", group_id, "Allow").
	   Find(&appPKeys)

	var users []models.User
	db.Model(&models.User{}).
	   Where("id IN (?)", user_ids).
	   Find(&users)

	
	for i := 0; i < len(users); i++ {
		for j := 0; j < len(appPKeys); j++ {
			if (users[i].AuthorizedGroupCount(appPKeys[j].Id, db) == 1) {
				cs_id := users[i].GetLatestActiveSession(db)
				if (cs_id != nil) {
					eventpublisher.PublishEvent(users[i].ID, cs_id, &appPKeys[j].Id, "AccessPolicyChanged", "group_deleted", mq)
				}
			}
		}
	}
}
