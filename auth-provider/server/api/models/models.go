package models

import (
	"auth-provider-server/api/utility"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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

func (u User) AuthorizedGroupCount(app_id datatypes.BinUUID, db *gorm.DB) int {
	groups := db.Model(&ApplicationGroupPolicy{}).Select("group_id").Where("application_id = ? AND effect = ?", app_id, "Allow")

	var userGroup []UserGroup
	db.Model(&UserGroup{}).Where("user_id = ? AND group_id IN (?)", u.ID, groups).Find(&userGroup)

	return len(userGroup)
}

func (u User) IsAuthorized(app_id datatypes.BinUUID, db *gorm.DB) bool {
	groups := db.Model(&ApplicationGroupPolicy{}).Select("group_id").Where("application_id = ? AND effect = ?", app_id, "Blocked")

	var userGroup []UserGroup
	db.Model(&UserGroup{}).Where("user_id = ? AND group_id IN (?)", u.ID, groups).Find(&userGroup)

	return (u.AuthorizedGroupCount(app_id, db) > 0) && (len(userGroup) == 0)
}

func (u User) GetLatestActiveSession(db *gorm.DB) *datatypes.BinUUID {
	var sessions []SSOSession
	db.Model(&SSOSession{}).Where("user_id = ? AND status = ?", u.ID, "Active").Find(&sessions)

	if len(sessions) <= 0 {
		return nil
	}

	return &(sessions[0].ID)
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
	ClientSecretHash      string            `json:"-"`
	Status                string            `json:"status"`
	LaunchUrl             string            `json:"launch_url"`
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

type AuditLog struct {
	ID            datatypes.BinUUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID()))"`
	EventType string
	ActorId *datatypes.BinUUID
	UserId *datatypes.BinUUID
	ApplicationId *datatypes.BinUUID
	SessionId *datatypes.BinUUID
	Result string
	Metadata datatypes.JSON
	IpAddress string
	CreatedAt time.Time
}

func AuditEvent(eventType, result string, user_id, app_id, session_id *datatypes.BinUUID, c *gin.Context, db *gorm.DB) {
	newUUID, _ := uuid.NewRandom()

	var actor_id *datatypes.BinUUID
	session_token, err := c.Cookie("ssid")
	if err != nil {
		actor_id = nil
	} else {
		var sessions []SSOSession

		db.Model(&SSOSession{}).
			Where("session_token_hash = ? AND status = ?", utility.HashToken(session_token), "Active").
			Find(&sessions)
		
		if len(sessions) == 0 {
			actor_id = nil
		} else {
			actor_id = &sessions[0].UserId
		}
	}

	audit_log := AuditLog{
		ID: datatypes.BinUUID(newUUID),
		EventType: eventType,
		Result: result,
		ActorId: actor_id,
		UserId: user_id,
		ApplicationId: app_id,
		SessionId: session_id,
		IpAddress: c.ClientIP(),
		CreatedAt: time.Now(),
	}

	db.Create(&audit_log)
}
