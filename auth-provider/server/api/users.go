package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

type PKeyPayload struct {
	Id datatypes.BinUUID `json:"id" binding:"required"`
}

type User struct {
	Id           datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	PasswordHash string            `json:"password_hash"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (payload CreateUserPayload) GetDBStruct() (User, error) {
	hashedPass, err := HashPassword(payload.Password)
	return User{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: hashedPass,
		Status:       "Inactive",
	}, err
}

func (payload UpdateUserPayload) GetDBStruct() (User, User, error) {
	hashedPass, err := HashPassword(payload.Password)
	return User{
			Id: payload.Id,
		},
		User{
			Name:         payload.Name,
			Email:        payload.Email,
			PasswordHash: hashedPass,
			Status:       payload.Status,
		}, err
}

func (payload PKeyPayload) GetUserModel() User {
	return User{Id: payload.Id}
}

func CreateUserRequest(c *gin.Context, db *gorm.DB) {
	var userPayload CreateUserPayload

	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}

	userModel, err := userPayload.GetDBStruct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}

	if err := db.Create(&userModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User successfully created",
	})
}

func GetAllUsersRequest(c *gin.Context, db *gorm.DB) {
	var users []User

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func UpdateUserRequest(c *gin.Context, db *gorm.DB) {
	var userPayload UpdateUserPayload

	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}

	userQuery, updatedUserModel, err := userPayload.GetDBStruct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}

	if err := db.Model(&userQuery).Updates(updatedUserModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User successfully updated",
	})
}

func DeleteUserRequest(c *gin.Context, db *gorm.DB) {
	var userPKey PKeyPayload

	if err := c.ShouldBindJSON(&userPKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete user: " + err.Error(),
		})
		return
	}

	userModel := userPKey.GetUserModel()
	if err := db.Delete(&userModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User successfully deleted",
	})
}
