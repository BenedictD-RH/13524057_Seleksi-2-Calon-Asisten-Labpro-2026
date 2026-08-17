package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID           datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
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
	ID          datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (Group) TableName() string {
	return "groups_"
}

type UserGroup struct {
	ID        datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	UserId    datatypes.BinUUID `json:"user_id"`
	GroupId   datatypes.BinUUID `json:"group_id"`
	CreatedAt time.Time         `json:"created_at"`
}

func (u User) GetGroupsUUID(db *gorm.DB) *[]datatypes.BinUUID {
	var groups []UserGroup
	db.Where("user_id = ?", u.ID).Find(&groups)
	groupsUUID := make([]datatypes.BinUUID, len(groups))
	for i := 0; i < len(groups); i++ {
		groupsUUID[i] = groups[i].GroupId
	}
	return &groupsUUID
}

type Application struct {
	ID                    datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name                  string            `json:"name"`
	ClientId              string            `json:"client_id"`
	ClientSecretHash      string            `json:"client_secret_hash"`
	Status                string            `json:"status"`
	LaunchUrl             string            `json:"redirect_url"`
	LogoutNotificationUrl string            `json:"logout_notification_url"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}



type AuthorizationCode struct {
	ID            datatypes.BinUUID `gorm:"default:UUID_TO_BIN(UUID())"`
	CodeHash      string
	CodeChallenge string
	UserId        datatypes.BinUUID
	ApplicationId datatypes.BinUUID
	SsoSessionId  datatypes.BinUUID
	RedirectUri   string
	CreatedAt     time.Time
	ExpiresAt     time.Time
	UsedAt        *time.Time
}

func (auth_code AuthorizationCode) useAuthCode(db *gorm.DB) {
	t := time.Now()
	db.Model(&AuthorizationCode{}).Where("id = ?", auth_code.ID).Updates(AuthorizationCode{UsedAt: &t})
}

func (auth_code AuthorizationCode) IsValid() bool {
	return auth_code.ExpiresAt.After(time.Now()) && (auth_code.UsedAt == nil)
}



type AccessToken struct {
	ID            datatypes.BinUUID `gorm:"default:UUID_TO_BIN(UUID())"`
	TokenHash     string
	UserId        datatypes.BinUUID
	ApplicationId datatypes.BinUUID
	SsoSessionId  datatypes.BinUUID
	Scopes        string
	Status        string
	CreatedAt     time.Time
	ExpiresAt     time.Time
	RevokedAt     *time.Time
}

func (token AccessToken) IsValid() bool {
	return token.Status == "Active" && token.ExpiresAt.After(time.Now()) && (token.RevokedAt == nil)
}

func (token AccessToken) UpdateStatus(db *gorm.DB) {
	if token.ExpiresAt.Before(time.Now()) {
		db.Model(&AccessToken{}).Where("id = ?", token.ID).Updates(AccessToken{Status: "Expired"})
	}
}

type SSOSession struct {
	ID               datatypes.BinUUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID()))"`
	UserId           datatypes.BinUUID
	SessionTokenHash string
	Status           string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	LastActivityAt   *time.Time
	RevokedAt        *time.Time
	RevokedReason    string
	IpAddress        string
	UserAgent        string
}

func (SSOSession) TableName() string {
	return "sso_sessions"
}

func (session SSOSession) IsValid() bool {
	return session.Status == "Active" && session.ExpiresAt.After(time.Now()) && (session.RevokedAt == nil)
}

func (session SSOSession) UpdateStatus(db *gorm.DB) {
	if session.ExpiresAt.Before(time.Now()) {
		db.Model(&SSOSession{}).Where("id = ?", session.ID).Updates(SSOSession{Status: "Expired"})
	}
}

func (session SSOSession) MarkActivity(db *gorm.DB) {
	t := time.Now()
	db.Model(&SSOSession{}).Where("id = ?", session.ID).Updates(SSOSession{LastActivityAt: &t})
}

type ApplicationRedirectURI struct {
	ID            datatypes.BinUUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID()))"`
	ApplicationId datatypes.BinUUID
	RedirectUri   string
	CreatedAt     time.Time
}

func (ApplicationRedirectURI) TableName() string {
	return "application_redirect_uris"
}

type ApplicationGroupPolicy struct {
	ID            datatypes.BinUUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID()))"`
	ApplicationId datatypes.BinUUID
	GroupId   	  datatypes.BinUUID
	Effect 		  string
	CreatedAt     time.Time
}

func (ApplicationGroupPolicy) TableName() string {
	return "application_group_policies"
}



