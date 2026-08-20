package controlpanel

import (
	eventpublisher "auth-provider-server/api/event-publisher"
	"auth-provider-server/api/models"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AppGroupPayload struct {
	AppId   datatypes.BinUUID `json:"app_id" binding:"required"`
	GroupId datatypes.BinUUID `json:"group_id" binding:"required"`
	Effect  string            `json:"effect"`
}

type GroupPolicyFromApp struct {
	GroupId     datatypes.BinUUID `json:"group_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Effect      string            `json:"effect"`
	CreatedAt   time.Time         `json:"created_at"`
}

type AppPolicyFromGroup struct {
	ApplicationId datatypes.BinUUID `json:"application_id"`
	Name          string            `json:"name"`
	ClientId      string            `json:"client_id"`
	Effect        string            `json:"effect"`
	CreatedAt     time.Time         `json:"created_at"`
}

func (payload AppGroupPayload) GetDBStructCreate() (models.ApplicationGroupPolicy, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.ApplicationGroupPolicy{}, err
	}
	return models.ApplicationGroupPolicy{
		ID:            datatypes.BinUUID(newUUID),
		ApplicationId: payload.AppId,
		GroupId:       payload.GroupId,
		Effect:        "Allow",
	}, nil
}

func AddAppGroupPolicyRequest(c *gin.Context, db *gorm.DB) {
	var appGroupPayload AppGroupPayload

	if err := c.ShouldBindJSON(&appGroupPayload); err != nil {
		models.AuditEvent("add_policy", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to add policy: " + err.Error(),
		})
		return
	}

	policyModel, err := appGroupPayload.GetDBStructCreate()
	if err != nil {
		models.AuditEvent("add_policy", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to add policy: " + err.Error(),
		})
		return
	}

	if err := db.Create(&policyModel).Error; err != nil {
		models.AuditEvent("add_policy", "failed", nil, &policyModel.ID, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to add policy: " + err.Error(),
		})
		return
	}

	models.AuditEvent("add_policy", "success", nil, &policyModel.ID, nil, c, db)
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Successfully added policy",
	})
}

func RemoveAppGroupPolicyRequest(c *gin.Context, db, mq *gorm.DB) {
	var appGroupPayload AppGroupPayload

	if err := c.ShouldBindJSON(&appGroupPayload); err != nil {
		models.AuditEvent("remove_policy", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove policy: " + err.Error(),
		})
		return
	}

	OnPolicyRemoved(appGroupPayload.AppId, appGroupPayload.GroupId, db, mq)

	db.Where("application_id = ? AND group_id = ?", appGroupPayload.AppId, appGroupPayload.GroupId).Delete(&models.ApplicationGroupPolicy{})

	models.AuditEvent("remove_policy", "success", nil, &appGroupPayload.AppId, nil, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully removed policy",
	})
}

func UpdateAppGroupPolicyRequest(c *gin.Context, db, mq *gorm.DB) {
	var appGroupPayload AppGroupPayload

	if err := c.ShouldBindJSON(&appGroupPayload); err != nil {
		models.AuditEvent("update_policy", "success", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update policy: " + err.Error(),
		})
		return
	}

	if appGroupPayload.Effect == "Blocked" {
		OnPolicyUpdateToBlocked(appGroupPayload.AppId, appGroupPayload.GroupId, db, mq)
	}

	db.Model(&models.ApplicationGroupPolicy{}).
		Where("application_id = ? AND group_id = ?", appGroupPayload.AppId, appGroupPayload.GroupId).
		Updates(models.ApplicationGroupPolicy{Effect: appGroupPayload.Effect})

	models.AuditEvent("update_policy", "success", nil, &appGroupPayload.AppId, nil, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully updated policy",
	})
}

func GetGroupPoliciesFromAppRequest(c *gin.Context, db *gorm.DB) {
	var policiesFromApp []GroupPolicyFromApp

	var appPKey PKeyPayload

	if err := c.ShouldBindJSON(&appPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to get groups: " + err.Error(),
		})
		return
	}

	db.Model(&models.ApplicationGroupPolicy{}).
		Select("application_group_policies.group_id, groups_.name, groups_.description, application_group_policies.effect, application_group_policies.created_at").
		Joins("LEFT JOIN groups_ ON groups_.id = application_group_policies.group_id").
		Where("application_id = ?", appPKey.Id).
		Find(&policiesFromApp)

	c.JSON(http.StatusOK, policiesFromApp)
}

func GetAppPoliciesFromGroupRequest(c *gin.Context, db *gorm.DB) {
	var policiesFromApp []AppPolicyFromGroup

	var groupPKey PKeyPayload

	if err := c.ShouldBindJSON(&groupPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to get apps: " + err.Error(),
		})
		return
	}

	db.Model(&models.ApplicationGroupPolicy{}).
		Select("application_group_policies.application_id , applications.name, applications.client_id, application_group_policies.effect, application_group_policies.created_at").
		Joins("LEFT JOIN applications ON applications.id = application_group_policies.application_id").
		Where("group_id = ?", groupPKey.Id).
		Find(&policiesFromApp)

	c.JSON(http.StatusOK, policiesFromApp)
}

func GetGroupsNotInAppPoliciesRequest(c *gin.Context, db *gorm.DB) {
	var groups []GroupCompact

	var appPKey PKeyPayload

	if err := c.ShouldBindJSON(&appPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	groupsInAppPolicy := db.Model(&models.Group{}).
		Select("groups_.id").
		Joins("RIGHT JOIN application_group_policies ON application_group_policies.group_id = groups_.id").
		Where("application_group_policies.application_id = ?", appPKey.Id)

	db.Model(&models.Group{}).Select("id, name, description").Where("id NOT IN (?)", groupsInAppPolicy).Find(&groups)

	c.JSON(http.StatusOK, groups)
}

func GetAppsNotInGroupPoliciesRequest(c *gin.Context, db *gorm.DB) {
	var apps []AppCompact

	var groupPKey PKeyPayload

	if err := c.ShouldBindJSON(&groupPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to remove user from group: " + err.Error(),
		})
		return
	}

	groupsInAppPolicy := db.Model(&models.Application{}).
		Select("applications.id").
		Joins("RIGHT JOIN application_group_policies ON application_group_policies.application_id = applications.id").
		Where("application_group_policies.group_id = ?", groupPKey.Id)

	db.Model(&models.Application{}).Select("id, name, client_id").Where("id NOT IN (?)", groupsInAppPolicy).Find(&apps)

	c.JSON(http.StatusOK, apps)
}

func GetGroupPoliciesFieldFromAppRequest(c *gin.Context) {
	groupPolicyFromApp := GroupPolicyFromApp{}

	t := reflect.TypeOf(groupPolicyFromApp)
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"fields": fieldNames,
	})
}

func GetAppPoliciesFieldFromGroupRequest(c *gin.Context) {
	appPolicyFromGroup := AppPolicyFromGroup{}

	t := reflect.TypeOf(appPolicyFromGroup)
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"fields": fieldNames,
	})
}


func OnPolicyUpdateToBlocked(app_id, group_id datatypes.BinUUID, db, mq *gorm.DB) {
	var groups []models.Group

	var policies []models.ApplicationGroupPolicy

	db.Model(&models.ApplicationGroupPolicy{}).
	   Where("application_id = ? AND group_id = ?", app_id, group_id).
	   Find(&policies)

	if len(policies) != 1 {
		return
	}

	if policies[0].Effect == "Blocked" {
		return
	}

	db.Model(&models.Group{}).Where("id = ?", policies[0].GroupId).Find(&groups)

	if len(groups) != 1 {
		return
	}

	var users []models.User

	user_ids := db.Model(&models.UserGroup{}).Select("user_id").Where("group_id = ?", groups[0].ID)

	db.Model(&models.User{}).Where("id IN (?)", user_ids).Find(&users)

	for i := 0; i < len(users); i++ {
		if (users[i].IsAuthorized(policies[0].ApplicationId, db)) {
			cs_id := users[i].GetLatestActiveSession(db)
			if (cs_id != nil) {
				eventpublisher.PublishEvent(users[i].ID, cs_id, &policies[0].ApplicationId, "AccessPolicyChanged", "policy_updated", mq)
			}
		}
	}
}

func OnPolicyRemoved(app_id, group_id datatypes.BinUUID, db, mq *gorm.DB) {
	var groups []models.Group

	var policies []models.ApplicationGroupPolicy

	db.Model(&models.ApplicationGroupPolicy{}).
	   Where("application_id = ? AND group_id = ?", app_id, group_id).
	   Find(&policies)

	if len(policies) != 1 {
		return
	}

	if policies[0].Effect == "Allowed" {
		return
	}

	db.Model(&models.Group{}).Where("id = ?", policies[0].GroupId).Find(&groups)

	if len(groups) != 1 {
		return
	}

	var users []models.User

	user_ids := db.Model(&models.UserGroup{}).Where("group_id = ?", groups[0].ID)

	db.Model(&models.User{}).Where("id IN (?)", user_ids).Find(&users)

	for i := 0; i < len(users); i++ {
		if (users[i].IsAuthorized(policies[0].ApplicationId, db)) {
			cs_id := users[i].GetLatestActiveSession(db)
			if (cs_id != nil) {
				eventpublisher.PublishEvent(users[i].ID, cs_id, &policies[0].ApplicationId, "AccessPolicyChanged", "policy_removed", mq)
			}
		}
	}
}