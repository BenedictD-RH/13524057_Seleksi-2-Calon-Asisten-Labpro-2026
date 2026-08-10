package models

import (
	"time"

	"gorm.io/datatypes"
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

type CodeVerifier struct {
	ID           datatypes.BinUUID `gorm:"default:UUID_TO_BIN(UUID())"`
	State        string
	CodeVerifier string
	Status       string
	CreatedAt    time.Time
	ExpiresAt    time.Time
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