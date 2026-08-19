package eventpublisher

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

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


type Event struct {
	ID datatypes.BinUUID `json:"id"`
	EventType string `json:"event_type"`
	UserId datatypes.BinUUID `json:"user_id"`
	CentralSessionId *datatypes.BinUUID `json:"central_session_id"`
	ApplicationId *datatypes.BinUUID `json:"application_id"`
	Payload datatypes.JSON
	Status string
	CreatedAt time.Time
	PublishedAt *time.Time
}


func PublishEvent(user_id datatypes.BinUUID, cs_id, app_id *datatypes.BinUUID, eventType, reason string, mq *gorm.DB) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return
	}

	eventPayload := EventPayload{
		EventId: datatypes.BinUUID(newUUID),
		EventType: eventType,
		UserId: user_id,
		CentralSessionId: cs_id,
		ApplicationId: app_id,
		Reason: reason,
		OccurredAt: time.Now(),
	}

	jsonData, _ := json.Marshal(eventPayload)
	
	t := time.Now()
	eventModel := Event{
		ID: datatypes.BinUUID(newUUID),
		EventType: eventType,
		UserId: user_id,
		CentralSessionId: cs_id,
		ApplicationId: app_id,
		Payload: jsonData,
		Status: "Published",
		CreatedAt: time.Now(),
		PublishedAt: &t,
	}

	mq.Create(&eventModel)
}