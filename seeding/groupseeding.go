package seeding

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Group struct {
	ID          datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}


type GroupData struct {
	Name string
	Description string
}

var SeedGroup = []GroupData {
	{"Administrators", "Has access to Control Panel"},
	{"App-A Users", "Is allowed access to App-A"},
	{"App-B Users", "Is allowed access to App-B"},
}

func SeedGroupsData(db *gorm.DB) {
	for i := 0; i < len(SeedGroup); i++ {
		newUUID, _ := uuid.NewRandom()

		group := Group{
			ID:           datatypes.BinUUID(newUUID),
			Name:         SeedGroup[0].Name,
			Description:        SeedGroup[0].Description,
		}

		db.Create(&group)
	}
}

type UserGroup struct {
	ID        datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	UserId    datatypes.BinUUID `json:"user_id"`
	GroupId   datatypes.BinUUID `json:"group_id"`
	CreatedAt time.Time         `json:"created_at"`
}

type UserGroupData struct {
	Email     string
	GroupName string
}
 
var SeedUserGroups = []UserGroupData{
	{"authadmin1@authprovider.com", "Administrators"},
	{"alice.johnson@example.com", "App-A Users"},
	{"bob.smith@example.com", "App-A Users"},
	{"bob.smith@example.com", "App-B Users"},
	{"carol.davis@example.com", "App-B Users"},
	{"eve.martinez@example.com", "App-A Users"},
	{"frank.thomas@example.com", "App-A Users"},
	{"grace.lee@example.com", "App-B Users"},
	{"henry.walker@example.com", "App-A Users"},
	{"henry.walker@example.com", "App-B Users"},
	{"ivy.robinson@example.com", "App-A Users"},
	{"jack.anderson@example.com", "App-B Users"},
	{"karen.clark@example.com", "App-A Users"},
	{"liam.rodriguez@example.com", "App-B Users"},
	{"mia.lewis@example.com", "App-A Users"},
	{"noah.walker@example.com", "App-A Users"},
	{"olivia.hall@example.com", "App-B Users"},
	{"olivia.hall@example.com", "App-A Users"},
	{"peter.young@example.com", "App-A Users"},
	{"rachel.king@example.com", "App-B Users"},
	{"sam.wright@example.com", "App-A Users"},
	{"tina.scott@example.com", "App-A Users"},
	{"tina.scott@example.com", "App-B Users"},
	{"uma.green@example.com", "App-B Users"},
	{"victor.adams@example.com", "App-A Users"},
	{"wendy.baker@example.com", "App-A Users"},
	{"xavier.nelson@example.com", "App-B Users"},
	{"yara.carter@example.com", "App-A Users"},
	{"zane.mitchell@example.com", "App-A Users"},
	{"amy.perez@example.com", "App-B Users"},
	{"chloe.turner@example.com", "App-A Users"},
	{"chloe.turner@example.com", "App-B Users"},
	{"dylan.phillips@example.com", "App-A Users"},
	{"ella.campbell@example.com", "App-A Users"},
	{"felix.parker@example.com", "App-B Users"},
	{"hugo.edwards@example.com", "App-A Users"},
	{"iris.collins@example.com", "App-A Users"},
	{"james.stewart@example.com", "App-B Users"},
	{"kylie.sanchez@example.com", "App-A Users"},
	{"molly.rogers@example.com", "App-A Users"},
	{"molly.rogers@example.com", "App-B Users"},
	{"nathan.reed@example.com", "App-A Users"},
}

func SeedUserGroupsData (db *gorm.DB) {
	for i:= 0; i < len(SeedUserGroups); i++ {
		
		var userPKey datatypes.BinUUID
		db.Model(&User{}).Where("email = ?", SeedUserGroups[0].Email).Pluck("id", &userPKey)

		var groupPKey datatypes.BinUUID
		db.Model(&Group{}).Where("name = ?", SeedUserGroups[0].GroupName).Pluck("id", &groupPKey)

		newUUID, _ := uuid.NewRandom()

		usergroup := UserGroup{
			ID: datatypes.BinUUID(newUUID),
			UserId: userPKey,
			GroupId: groupPKey,
		}

		db.Create(&usergroup)
	}
}