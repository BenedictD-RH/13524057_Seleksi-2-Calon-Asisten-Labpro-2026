package seeding

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func main() {
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

	SeedUsersData(db)
	fmt.Println("Seeded user data")
	SeedGroupsData(db)
	fmt.Println("Seeded group data")
	SeedUserGroupsData(db)
	fmt.Println("Seeded user_group data")
	SeedApplicationsData(db)
	fmt.Println("Seeded application data")
	SeedPoliciesData(db)
	fmt.Println("Seeded application_group_policy data")
}