package main

import (
	"app-a-backend/api/app"
	"app-a-backend/api/auth"
	"app-a-backend/api/models"
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
  db_pass := os.Getenv("APP_B_DB_PASS")
  dsn := fmt.Sprintf("root:%s@tcp(localhost:8790)/app-db?charset=utf8mb4&parseTime=True&loc=Local", db_pass)
  db, db_err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if (db_err != nil) {
    fmt.Println(db_err);
    return;
  }
  if (db == nil) {
    return;
  } 

  r := gin.Default()
  
  r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("AUTH_SERVER_URL"), os.Getenv("AUTH_PORTAL_URL"), os.Getenv("FRONTEND_URI")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

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

  r.GET("/session", func(c *gin.Context) {
    auth.GetLocalSessionRequest(c, db)
  })

  r.POST("/logout", func(c *gin.Context) {
    auth.LogoutRequest(c, db)
  })

  r.POST("/internal/logout", func(c *gin.Context) {
    auth.BackChannelLogoutRequest(c, db)
  })

  r.GET("/logs/activity", func(c *gin.Context) {
    models.GetAllActivityLogs(c, db)
  })

  r.GET("/logs/event", func(c *gin.Context) {
    models.GetAllProcessedEventLogs(c, db)
  })

  r.RunTLS(":8791", "./certs/localhost.pem", "./certs/localhost-key.pem")
}