package models

import (
	"time"

	"gorm.io/datatypes"
)


type LocalSession struct {
	Id datatypes.BinUUID 		`gorm:"default:UUID_TO_BIN(UUID())"`
	SessionTokenHash string
	ExternalUserId datatypes.BinUUID
	CentralSessionId datatypes.BinUUID
	ApplicationId *datatypes.BinUUID
	Status string
	CreatedAt time.Time
	ExpiresAt time.Time
	LastActivityAt *time.Time
	RevokedAt *time.Time
	RevokedReason string
}