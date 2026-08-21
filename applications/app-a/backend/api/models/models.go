package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LocalSession struct {
	ID               datatypes.BinUUID `gorm:"default:UUID_TO_BIN(UUID())"`
	SessionTokenHash string
	ExternalUserId   datatypes.BinUUID
	CentralSessionId datatypes.BinUUID
	ApplicationId    *datatypes.BinUUID
	Status           string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	LastActivityAt   *time.Time
	RevokedAt        *time.Time
	RevokedReason    string
}

func (session LocalSession) IsValid() bool {
	return session.Status == "Active" && session.ExpiresAt.After(time.Now()) && session.RevokedAt == nil
}

func (session LocalSession) UpdateStatus(db *gorm.DB) {
	if session.ExpiresAt.Before(time.Now()) {
		db.Model(&LocalSession{}).Where("id = ?", session.ID).Updates(LocalSession{Status: "Expired"})
	}
}

func (session LocalSession) MarkActivity(db *gorm.DB) {
	t := time.Now()
	db.Model(&LocalSession{}).Where("id = ?", session.ID).Updates(LocalSession{LastActivityAt: &t})
}

type CodeVerifier struct {
	ID           datatypes.BinUUID `gorm:"default:UUID_TO_BIN(UUID())"`
	State        string
	CodeVerifier string
	Status       string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func (code CodeVerifier) IsValid() bool {
	return code.ExpiresAt.After(time.Now())
}

type ProfileCache struct {
	ExternalUserId   datatypes.BinUUID
	Name string
	Email string
	GroupsList datatypes.JSON
	SyncedAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (ProfileCache) TableName() string {
	return "profile_cache"
}

type ProcessedEvent struct {
	EventID datatypes.BinUUID
	EventType string
	ProcessedAt time.Time
	Result string
}

func (a ProcessedEvent) ToString() string {
	return fmt.Sprintf("[%s][%s] at %s", 
		a.EventType, 
		a.Result,
		a.ProcessedAt)
}

type ActivityLog struct {
	ID datatypes.BinUUID
	EventType string
	ExternalUserId *datatypes.BinUUID
	CentralSessionId *datatypes.BinUUID
	Result string
	CreatedAt time.Time
}

func LogActivity(eventType, result string, external_user_id, central_session_id *datatypes.BinUUID, db *gorm.DB) {
	newUUID, _ := uuid.NewRandom()


	activity_log := ActivityLog{
		ID: datatypes.BinUUID(newUUID),
		EventType: eventType,
		ExternalUserId: external_user_id,
		CentralSessionId: central_session_id,
		Result: result,
		CreatedAt: time.Now(),
	}

	db.Create(&activity_log)
}

func (a ActivityLog) ToString() string {
	var userString, sessionString string
	if a.ExternalUserId != nil {
		userString = "User: " + a.ExternalUserId.String()
		if a.CentralSessionId != nil {
			userString +=", "
		}
	}
	if a.CentralSessionId != nil {
		sessionString = "Session: " + a.CentralSessionId.String()
	}

	return fmt.Sprintf("[%s][%s] by (%s%s) at %s", 
		a.EventType, 
		a.Result,
		userString,
		sessionString,
		a.CreatedAt)
}

func GetAllActivityLogs(c *gin.Context, db *gorm.DB) {
	var logs []ActivityLog
	db.Model(&ActivityLog{}).Order("created_at DESC").Find(&logs)

	var stringifiedLogs []string
	for i := 0; i < len(logs);i++ {
		stringifiedLogs = append(stringifiedLogs, logs[i].ToString())
	}

	c.JSON(http.StatusOK, gin.H{
		"log" : stringifiedLogs,
	})
}

func GetAllProcessedEventLogs(c *gin.Context, db *gorm.DB) {
	var logs []ProcessedEvent
	db.Model(&ProcessedEvent{}).Order("processed_at DESC").Find(&logs)

	var stringifiedLogs []string
	for i := 0; i < len(logs);i++ {
		stringifiedLogs = append(stringifiedLogs, logs[i].ToString())
	}

	c.JSON(http.StatusOK, gin.H{
		"log" : stringifiedLogs,
	})
}