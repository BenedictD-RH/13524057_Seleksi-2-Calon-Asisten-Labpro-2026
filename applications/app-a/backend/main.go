package main

import (
	"fmt"
	"os"
  "app-a-backend/api/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  "github.com/joho/godotenv"
)



func main() {
  // Remove for docker
  err := godotenv.Load("../.env");
  if (err != nil) {
    fmt.Println("Missing .env file: " + err.Error())
    return
  }
  // Remove for docker
  db_pass := os.Getenv("APP_A_DB_PASS")
  dsn := fmt.Sprintf("root:%s@tcp(localhost:8690)/app-db?charset=utf8mb4&parseTime=True&loc=Local", db_pass)
  db, db_err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if (db_err != nil) {
    fmt.Println(db_err);
    return;
  }
  if (db == nil) {
    return;
  } 

  r := gin.Default()

  r.GET("/login", func(c *gin.Context) {
    auth.LoginRequest(c, db)
  })

  r.Run()
}