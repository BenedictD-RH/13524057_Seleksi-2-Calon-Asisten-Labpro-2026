package helper

import (
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

type ApplicationRedirectURI struct {
	ID            datatypes.BinUUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID()))"`
	ApplicationId datatypes.BinUUID
	RedirectUri   string
	CreatedAt     time.Time
}

func SeedApplicationsData(db *gorm.DB) {
	newUUID, _ := uuid.NewRandom()

	client_a_secret_hash, _ := HashPassword(os.Getenv("CLIENT_A_SECRET"))

	app_a := Application{
		ID: datatypes.BinUUID(newUUID),
		Name: "App-A",
		ClientId: os.Getenv("CLIENT_A_ID"),
		ClientSecretHash: client_a_secret_hash,
		Status: "Active",
		LaunchUrl: os.Getenv("APP_A_FRONTEND"),
		LogoutNotificationUrl: os.Getenv("APP_A_BACKEND_2") + "/internal/logout",
	}

	newUUID, _ = uuid.NewRandom()
	app_a_redirect_uri := ApplicationRedirectURI{
		ID: datatypes.BinUUID(newUUID),
		ApplicationId: app_a.ID,
		RedirectUri: os.Getenv("APP_A_BACKEND") + "/auth/callback",
	}

	db.Create(&app_a)
	db.Create(&app_a_redirect_uri)

	newUUID, _ = uuid.NewRandom()

	client_b_secret_hash, _ := HashPassword(os.Getenv("CLIENT_B_SECRET"))

	app_b := Application{
		ID: datatypes.BinUUID(newUUID),
		Name: "App-B",
		ClientId: os.Getenv("CLIENT_B_ID"),
		ClientSecretHash: client_b_secret_hash,
		Status: "Active",
		LaunchUrl: os.Getenv("APP_B_FRONTEND"),
		LogoutNotificationUrl: os.Getenv("APP_B_BACKEND_2") + "/internal/logout",
	}

	newUUID, _ = uuid.NewRandom()
	app_b_redirect_uri := ApplicationRedirectURI{
		ID: datatypes.BinUUID(newUUID),
		ApplicationId: app_b.ID,
		RedirectUri: os.Getenv("APP_B_BACKEND") + "/auth/callback",
	}

	db.Create(&app_b)
	db.Create(&app_b_redirect_uri)
}

type ApplicationGroupPolicy struct {
	ID            datatypes.BinUUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID()))"`
	ApplicationId datatypes.BinUUID
	GroupId   	  datatypes.BinUUID
	Effect 		  string
	CreatedAt     time.Time
}

type PolicyData struct {
	ClientID string
	GroupName string
	Effect string
}

var SeedPolicies = []PolicyData{
	{"CLIENT_A_ID", "App-A Users", "Allow"},
	{"CLIENT_B_ID", "App-B Users", "Allow"},
	{"CLIENT_A_ID", "Administrators", "Allow"},
	{"CLIENT_B_ID", "Administrators", "Allow"},
}

func SeedPoliciesData(db *gorm.DB) {
	for i := 0; i < len(SeedPolicies); i++ {
		newUUID, _ := uuid.NewRandom()

		var appPKey Application
		db.Model(&Application{}).Where("client_id = ?", os.Getenv(SeedPolicies[i].ClientID)).First(&appPKey)

		var groupPKey Group
		db.Model(&Group{}).Where("name = ?", SeedPolicies[i].GroupName).First(&groupPKey)

		policy := ApplicationGroupPolicy{
			ID: datatypes.BinUUID(newUUID),
			ApplicationId: appPKey.ID,
			GroupId: groupPKey.ID,
			Effect: SeedPolicies[i].Effect,
		}

		db.Create(&policy)
	}
}