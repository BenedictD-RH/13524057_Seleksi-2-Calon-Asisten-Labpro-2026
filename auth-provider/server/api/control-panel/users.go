package controlpanel

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"

	eventpublisher "auth-provider-server/api/event-publisher"
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CreateUserPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8"`
}

type UpdateUserPayload struct {
	Id       datatypes.BinUUID `json:"id" binding:"required"`
	Name     string            `json:"name"`
	Email    string            `json:"email" binding:"omitempty,email"`
	Password string            `json:"password" binding:"omitempty,gte=8"`
	Status   string            `json:"status"`
}

type UserCompact struct {
	Id       datatypes.BinUUID `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
}

type PKeyPayload struct {
	Id datatypes.BinUUID `json:"id" binding:"required"`
}

func (payload CreateUserPayload) GetDBStruct() (models.User, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.User{}, err
	}
	hashedPass, err := utility.HashPassword(payload.Password)
	return models.User{
		ID:           datatypes.BinUUID(newUUID),
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: hashedPass,
		Status:       "Active",
	}, err
}

func (payload UpdateUserPayload) GetDBStruct() (models.User, models.User, error) {
	var hashedPass string = ""
	var err error = nil
	if payload.Password != "" {
		hashedPass, err = utility.HashPassword(payload.Password)
	}

	return models.User{
			ID: payload.Id,
		},
		models.User{
			Name:         payload.Name,
			Email:        payload.Email,
			PasswordHash: hashedPass,
			Status:       payload.Status,
		}, err
}

func (payload PKeyPayload) GetUserModel() models.User {
	return models.User{ID: payload.Id}
}

func CreateUserRequest(c *gin.Context, db *gorm.DB) {
	var userPayload CreateUserPayload

	if err := c.ShouldBindJSON(&userPayload); err != nil {
		models.AuditEvent("create_user", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}
	userModel, err := userPayload.GetDBStruct()
	if err != nil {
		models.AuditEvent("create_user", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}

	if err := db.Create(&userModel).Error; err != nil {
		models.AuditEvent("create_user", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}

	models.AuditEvent("create_user", "success", &userModel.ID, nil, nil, c, db)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User successfully created",
	})
}

func GetAllUsersRequest(c *gin.Context, db *gorm.DB) {
	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserFields(c *gin.Context) {
	user := models.User{}

	t:= reflect.TypeOf(user)
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"fields" : fieldNames,
	})
}

func UpdateUserRequest(c *gin.Context, db, mq *gorm.DB) {
	var userPayload UpdateUserPayload

	if err := c.ShouldBindJSON(&userPayload); err != nil {
		models.AuditEvent("update_user", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}
	
	if (userPayload.Status == "Inactive") {
		OnSetUserInactiveUpdate(userPayload.Id, db, mq)
	}
	
	updateType := "update_user"
	userQuery, updatedUserModel, err := userPayload.GetDBStruct()
	if (updatedUserModel.PasswordHash != "") {
		updateType = "password_change"
		OnUserPasswordChangeUpdate(userPayload.Id, db, mq)
	}
	if err != nil {
		models.AuditEvent(updateType, "failed", &userPayload.Id, nil, nil, c, db)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}

	if err := db.Model(&userQuery).Updates(updatedUserModel).Error; err != nil {
		models.AuditEvent(updateType, "failed", &userPayload.Id, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}
	models.AuditEvent(updateType, "success", &userPayload.Id, nil, nil, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User successfully updated",
	})
}

func DeleteUserRequest(c *gin.Context, db, mq *gorm.DB) {
	var userPKey PKeyPayload

	if err := c.ShouldBindJSON(&userPKey); err != nil {
		models.AuditEvent("delete_user", "failed", nil, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete user: " + err.Error(),
		})
		return
	}

	OnDeleteUser(userPKey.Id, db, mq)

	userModel := userPKey.GetUserModel()
	if err := db.Delete(&userModel).Error; err != nil {
		models.AuditEvent("delete_user", "failed", &userPKey.Id, nil, nil, c, db)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete user: " + err.Error(),
		})
		return
	}
	models.AuditEvent("delete_user", "success", &userPKey.Id, nil, nil, c, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User successfully deleted",
	})
}


func OnUserPasswordChangeUpdate(user_id datatypes.BinUUID, db, mq *gorm.DB) {
	user := models.User{ID: user_id}
	cs_id := user.GetLatestActiveSession(db)
	if cs_id != nil {
		eventpublisher.PublishEvent(user.ID, cs_id, nil, "PasswordChanged", "user_password_changed", mq)
	}
	db.Model(&models.SSOSession{}).
	   Where("id = ?", cs_id).
	   Updates(models.SSOSession{Status: "Revoked", RevokedReason: "user_logout"})
}

func OnSetUserInactiveUpdate(user_id datatypes.BinUUID, db, mq *gorm.DB) {
	user := models.User{ID: user_id}
	cs_id := user.GetLatestActiveSession(db)
	if cs_id != nil {
		eventpublisher.PublishEvent(user.ID, cs_id, nil, "SessionRevoked", "user_set_inactive", mq)
	}
	db.Model(&models.SSOSession{}).
	   Where("id = ?", cs_id).
	   Updates(models.SSOSession{Status: "Revoked", RevokedReason: "user_logout"})
}

func OnDeleteUser(user_id datatypes.BinUUID, db, mq *gorm.DB) {
	user := models.User{ID: user_id}
	cs_id := user.GetLatestActiveSession(db)
	if cs_id != nil {
		eventpublisher.PublishEvent(user.ID, cs_id, nil, "SessionRevoked", "user_deleted", mq)
	}
	db.Model(&models.SSOSession{}).
	   Where("id = ?", cs_id).
	   Updates(models.SSOSession{Status: "Revoked", RevokedReason: "user_logout"})
}