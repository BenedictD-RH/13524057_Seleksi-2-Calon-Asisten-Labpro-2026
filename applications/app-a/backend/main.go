package main

import (
	"app-a-backend/api/app"
	"app-a-backend/api/auth"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
  r.Use(cors.Default())

  r.GET("/login", func(c *gin.Context) {
    auth.LoginRequest(c, db)
  })

  r.GET("/auth/callback", func(c *gin.Context) {
    auth.AuthCallbackRequest(c, db)
  })

  r.POST("/auth/register", func(c *gin.Context) {
    auth.RegisterAppToAuthProviderRequest(c, db)
  })

  r.GET("/users", func(c *gin.Context) {
    app.GetUserDataRequest(c, db)
  })

  r.Run(":8691")
}