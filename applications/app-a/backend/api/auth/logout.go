package auth

import (
	"app-a-backend/api/models"
	"app-a-backend/api/utility"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EventPayload struct {
	EventId datatypes.BinUUID `json:"event_id"`
	EventType string `json:"event_type"`
	UserId datatypes.BinUUID `json:"user_id"`
	CentralSessionId *datatypes.BinUUID `json:"central_session_id"`
	ApplicationId *datatypes.BinUUID `json:"application_id"`
	Reason string `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
	Metadata datatypes.JSON `json:"metadata"`
}

func LogoutRequest(c *gin.Context, db *gorm.DB) {
	local_session_token, err := c.Cookie("local_ssid")
	if err != nil {
		UnauthorizedResponse(c, "No local session found.")
		return
	}

	var localSession models.LocalSession
	db.Model(&models.LocalSession{}).
	   Where("session_token_hash = ?", utility.HashToken(local_session_token)).
	   First(&localSession)

	db.Model(&models.LocalSession{}).
	   Where("session_token_hash = ?", utility.HashToken(local_session_token)).
	   Updates(models.LocalSession{Status: "Revoked", RevokedReason: "user_logout"})

	c.SetCookie("local_ssid", "", -1, "/", "", false, true)

	models.LogActivity("user_logout", "sucess", &localSession.ExternalUserId, &localSession.CentralSessionId, db)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout successful",
	})
}

func BackChannelLogoutRequest(c *gin.Context, db *gorm.DB) {
	var eventPayload EventPayload

	if err := c.ShouldBindJSON(&eventPayload); err != nil {
		InternalServerErrorResponse(c)
		return
	}

	var processed_events []models.ProcessedEvent
	db.Model(&models.ProcessedEvent{}).Where("event_id = ?", eventPayload.EventId).Find(&processed_events)

	if len(processed_events) == 1 {
		if processed_events[0].Result == "Success" {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Logout successful",
			})
			return
		}
	}

	err := db.Model(&models.LocalSession{}).
			Where("external_user_id = ? AND central_session_id = ?", eventPayload.UserId, eventPayload.CentralSessionId).
			Updates(models.LocalSession{Status: "Revoked", RevokedReason: eventPayload.Reason}).Error
	
	result := "success"
	if err != nil {
		result = "failed"
	}

	if len(processed_events) < 1 {
		db.Create(&models.ProcessedEvent{
			EventID: eventPayload.EventId, 
			EventType: eventPayload.EventType, 
			ProcessedAt: time.Now(),
			Result: result,
		})
	} else {
		db.Model(&models.ProcessedEvent{}).
		Where("event_id = ?", eventPayload.EventId).
		Updates(models.ProcessedEvent{
			ProcessedAt: time.Now(),
			Result: result,
		})
	}

	models.LogActivity("backchannel_logout", result, &eventPayload.UserId, eventPayload.CentralSessionId, db)

	c.JSON(http.StatusOK, gin.H{
		"status":  result,
		"message": "Logout " + result,
	})
}