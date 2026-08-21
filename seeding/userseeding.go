package seeding

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserData struct {
	Name     string
	Email    string
	Password string
	Status   string
}

type User struct {
	ID           datatypes.BinUUID `json:"id" gorm:"default:UUID_TO_BIN(UUID())"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	PasswordHash string            `json:"password_hash"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

var SeedUsers = []UserData{
	{"Admin1", "authadmin1@authprovider.com", "password123", "Active"},
	{"Alice Johnson", "alice.johnson@example.com", "correcthorse1", "Active"},
	{"Bob Smith", "bob.smith@example.com", "sunsetbeach22", "Active"},
	{"Carol Davis", "carol.davis@example.com", "bluemountain7", "Active"},
	{"David Wilson", "david.wilson@example.com", "quietriver99", "Inactive"},
	{"Eve Martinez", "eve.martinez@example.com", "goldenfalcon3", "Active"},
	{"Frank Thomas", "frank.thomas@example.com", "silverwolf88", "Active"},
	{"Grace Lee", "grace.lee@example.com", "crimsonrose15", "Inactive"},
	{"Henry Walker", "henry.walker@example.com", "northwind44", "Active"},
	{"Ivy Robinson", "ivy.robinson@example.com", "brightstar21", "Active"},
	{"Jack Anderson", "jack.anderson@example.com", "stonebridge6", "Active"},
	{"Karen Clark", "karen.clark@example.com", "purplehaze77", "Active"},
	{"Liam Rodriguez", "liam.rodriguez@example.com", "windyharbor12", "Active"},
	{"Mia Lewis", "mia.lewis@example.com", "autumnleaf55", "Inactive"},
	{"Noah Walker", "noah.walker@example.com", "deepocean33", "Active"},
	{"Olivia Hall", "olivia.hall@example.com", "lucky7cat", "Active"},
	{"Peter Young", "peter.young@example.com", "greenvalley9", "Active"},
	{"Quinn Allen", "quinn.allen@example.com", "midnightowl2", "Inactive"},
	{"Rachel King", "rachel.king@example.com", "coralreef66", "Active"},
	{"Sam Wright", "sam.wright@example.com", "shadowfox18", "Active"},
	{"Tina Scott", "tina.scott@example.com", "winterpine41", "Active"},
	{"Uma Green", "uma.green@example.com", "amberglow5", "Active"},
	{"Victor Adams", "victor.adams@example.com", "thunderbolt8", "Inactive"},
	{"Wendy Baker", "wendy.baker@example.com", "mossygrove14", "Active"},
	{"Xavier Nelson", "xavier.nelson@example.com", "scarletmoon27", "Active"},
	{"Yara Carter", "yara.carter@example.com", "frostyoak31", "Active"},
	{"Zane Mitchell", "zane.mitchell@example.com", "ironclad19", "Active"},
	{"Amy Perez", "amy.perez@example.com", "velvetsky24", "Active"},
	{"Brian Roberts", "brian.roberts@example.com", "rustyanchor10", "Inactive"},
	{"Chloe Turner", "chloe.turner@example.com", "lunarbeam36", "Active"},
	{"Dylan Phillips", "dylan.phillips@example.com", "cedarpath42", "Active"},
	{"Ella Campbell", "ella.campbell@example.com", "honeybee29", "Active"},
	{"Felix Parker", "felix.parker@example.com", "graniteedge13", "Active"},
	{"Gina Evans", "gina.evans@example.com", "violetdusk50", "Inactive"},
	{"Hugo Edwards", "hugo.edwards@example.com", "maplewood17", "Active"},
	{"Iris Collins", "iris.collins@example.com", "pebblebeach4", "Active"},
	{"James Stewart", "james.stewart@example.com", "falconcrest35", "Active"},
	{"Kylie Sanchez", "kylie.sanchez@example.com", "driftwood23", "Active"},
	{"Leo Morris", "leo.morris@example.com", "emberglow16", "Inactive"},
	{"Molly Rogers", "molly.rogers@example.com", "willowbend38", "Active"},
	{"Nathan Reed", "nathan.reed@example.com", "canyonwind20", "Active"},
}

func SeedUsersData(db *gorm.DB) {
	for i := 0; i < len(SeedUsers); i++ {
		newUUID, _ := uuid.NewRandom()

		hashedPass, _ := HashPassword(SeedUsers[0].Password)

		user := User{
			ID:           datatypes.BinUUID(newUUID),
			Name:         SeedUsers[0].Name,
			Email:        SeedUsers[0].Email,
			PasswordHash: hashedPass,
			Status:       SeedUsers[0].Status,
		}

		db.Create(&user)
	}
}
