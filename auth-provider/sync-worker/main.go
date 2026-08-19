package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/joho/godotenv"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

type EventDelivery struct {
	ID datatypes.BinUUID `json:"id"`
	EventId datatypes.BinUUID `json:"event_type"`
	ApplicationId datatypes.BinUUID `json:"application_id"`
	Status string
	AttemptCount int
    LastAttemptAt *time.Time
    NextRetryAt *time.Time
    ProcessedAt *time.Time
    LastError string
}

func (EventDelivery) TableName() string {
	return "event_deliveries"
}

func CreateEventDelivery(event Event, app_id datatypes.BinUUID) (*EventDelivery) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil
	}

	return &EventDelivery{
		ID: datatypes.BinUUID(newUUID),
		EventId: event.ID,
		ApplicationId: app_id,
		Status: "Processing",
		AttemptCount: 0,
	}
}

func EventDeliverer(ctx context.Context, db, mq *gorm.DB) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	fmt.Println("Event Deliverer running...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Event Deliverer shutting down gracefully...")
			return
		case <-ticker.C:
			var published_events []Event
			mq.Model(&Event{}).Where("status = ?", "Published").Find(&published_events)
			
			for i:= 0; i < len(published_events); i++ {
				if (published_events[0].ApplicationId == nil) {
					var app_id []datatypes.BinUUID
					db.Table("applications").Pluck("id", &app_id)

					for j := 0; j < len(app_id); j++ {
						mq.Create(CreateEventDelivery(published_events[i], app_id[j]))
					}
				} else {
					mq.Create(CreateEventDelivery(published_events[i], *published_events[i].ApplicationId))
				}

				mq.Model(&Event{}).Where("id = ?", published_events[i].ID).Update("status", "Processing")
			}
		}
	}
}

func EventProcessor(ctx context.Context, db, mq *gorm.DB) {
	delivery_interval := 5 * time.Second
	ticker := time.NewTicker(delivery_interval)
	defer ticker.Stop()
	max_attempt_count := 200


	fmt.Println("Event Processor running...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Event Processor shutting down gracefully...")
			return
		case <-ticker.C:
			var event_deliveries []EventDelivery

			mq.Model(&EventDelivery{}).Where("status = ?", "Processing").Find(&event_deliveries)

			for i:= 0; i < len(event_deliveries); i++ {
				var event Event
				mq.Model(&Event{}).Where("id = ?", event_deliveries[i].EventId).First(&event)

				var app_endpoint string

				db.Table("applications").
					Where("id = ?", event_deliveries[i].ApplicationId).
					Select("logout_notification_url").
					Row().Scan(&app_endpoint)

				if (event_deliveries[i].LastAttemptAt == nil) {
					mq.Model(&EventDelivery{}).Where("id = ?", event_deliveries[i].ID).Update("last_attempt_at", time.Now())
				}
				
				mq.Model(&EventDelivery{}).Where("id = ?", event_deliveries[i].ID).Update("attempt_count", event_deliveries[i].AttemptCount + 1)

				if (event_deliveries[i].AttemptCount < max_attempt_count) {
					req, _ := http.NewRequest("POST", app_endpoint, bytes.NewReader(event.Payload))

					req.Header.Set("Content-Type", "application/json")
					client := &http.Client{Timeout: 10 * time.Second}
					resp, err := client.Do(req)

					if (resp.StatusCode == http.StatusOK) {
						mq.Model(&EventDelivery{}).
						   Where("id = ?", event.ID).
						   Update("status", "Delivered")
						
						t := time.Now()

						mq.Model(&EventDelivery{}).
						   Where("id = ?", event_deliveries[i].ID).
						   Updates(EventDelivery{Status: "Success", ProcessedAt: &t})
					} else {
						t := time.Now().Add(delivery_interval)
						mq.Model(&EventDelivery{}).
						   Where("id = ?", event_deliveries[i].ID).
						   Updates(EventDelivery{NextRetryAt: &t, LastError: err.Error()})
					}
				} else {
					mq.Model(&EventDelivery{}).
						Where("id = ?", event_deliveries[i].ID).
						Updates(EventDelivery{Status: "MaxAttemptReached"})
				}
			}
		}
	}
}

func main() {
	// Remove for docker
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Missing .env file: " + err.Error())
		return
	}
	// Remove for docker
	db_pass := os.Getenv("AUTH_PROVIDER_DB_PASS")
	dsn := fmt.Sprintf("root:%s@tcp(localhost:5342)/auth-provider-db?charset=utf8mb4&parseTime=True&loc=Local", db_pass)
	db, db_err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if db_err != nil {
		fmt.Println(db_err)
		return
	}
	if db == nil {
		return
	}

	mq_pass := os.Getenv("MESSAGE_QUEUE_DB_PASS")
	dsn_mq := fmt.Sprintf("root:%s@tcp(localhost:5343)/queue?charset=utf8mb4&parseTime=True&loc=Local", mq_pass)
	mq, db_err := gorm.Open(mysql.Open(dsn_mq), &gorm.Config{})
	if db_err != nil {
		fmt.Println(db_err)
		return
	}
	if mq == nil {
		return
	}

	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go EventDeliverer(ctx, db, mq)

	go EventProcessor(ctx, db, mq)


	select {} 
	
}


