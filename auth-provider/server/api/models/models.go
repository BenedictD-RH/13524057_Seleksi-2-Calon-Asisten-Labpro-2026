package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	Id           datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	PasswordHash string            `json:"password_hash"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func (u User) IsActive() bool {
	return u.Status == "Active"
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

type UserGroup struct {
	Id datatypes.BinUUID 		`json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	UserId datatypes.BinUUID	`json:"user_id"`
	GroupId datatypes.BinUUID	`json:"group_id"`
	CreatedAt time.Time         `json:"created_at"`
}

type Application struct {
	Id datatypes.BinUUID		`json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name string          		`json:"name"`
	ClientId string				`json:"client_id"`
	ClientSecretHash string		`json:"client_secret_hash"`
	Status string            	`json:"status"`
	RedirectUrl string			`json:"redirect_url"`
	LogoutNotificationUrl string`json:"logout_notification_url"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type AuthorizationCode struct {
	Id 	datatypes.BinUUID 		`gorm:"default:UUID_TO_BIN(UUID())"`
	CodeHash string
	CodeChallenge string
	UserId datatypes.BinUUID 
	ApplicationId datatypes.BinUUID 
	SsoSessionId datatypes.BinUUID 
	RedirectUri string
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt *time.Time
}

type SSOSession struct {
	Id datatypes.BinUUID 		`gorm:"default:UUID_TO_BIN(UUID())"`
	UserId datatypes.BinUUID
	SessionTokenHash string
	Status string
	CreatedAt time.Time
	ExpiresAt time.Time
	LastActivityAt *time.Time
	RevokedAt *time.Time
	RevokedReason string
	IpAddress string
	UserAgent string
}

func (SSOSession) TableName() string {
	return "sso_sessions"
}