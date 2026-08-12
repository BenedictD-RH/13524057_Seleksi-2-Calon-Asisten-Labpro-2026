package models

import (
	"time"

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