package main

import (
	"fmt"
	"os"

	centralsessionserver "auth-provider-server/api/central-session-server"
	"auth-provider-server/api/control-panel"

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
  db_pass := os.Getenv("AUTH_PROVIDER_DB_PASS")
  dsn := fmt.Sprintf("root:%s@tcp(localhost:5342)/auth-provider-db?charset=utf8mb4&parseTime=True&loc=Local", db_pass)
  db, db_err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if (db_err != nil) {
    fmt.Println(db_err);
    return;
  }
  if (db == nil) {
    return;
  } 

  r := gin.Default()

  r.POST("/users", func(c *gin.Context) {
    controlpanel.CreateUserRequest(c, db)
  })

  r.GET("/users", func(c *gin.Context) {
    controlpanel.GetAllUsersRequest(c, db)
  })

  r.PATCH("/users", func(c *gin.Context) {
    controlpanel.UpdateUserRequest(c, db)
  })

  r.DELETE("/users", func(c *gin.Context) {
    controlpanel.DeleteUserRequest(c, db)
  })

  r.POST("/groups", func(c *gin.Context) {
    controlpanel.CreateGroupRequest(c, db)
  })

  r.GET("/groups", func(c *gin.Context) {
    controlpanel.GetAllGroupsRequest(c, db)
  })

  r.PATCH("/groups", func(c *gin.Context) {
    controlpanel.UpdateGroupRequest(c, db)
  })

  r.DELETE("/groups", func(c *gin.Context) {
    controlpanel.DeleteGroupRequest(c, db)
  })

  r.POST("/usergroups", func(c *gin.Context) {
    controlpanel.AddUserToGroup(c, db)
  })

  r.DELETE("/usergroups", func(c *gin.Context) {
    controlpanel.RemoveUserFromGroup(c, db)
  })

  r.POST("/login", func(c *gin.Context) {
    centralsessionserver.LoginRequest(c, db)
  })

  r.GET("/authorize", func(c *gin.Context){
    centralsessionserver.AuthorizeRequest(c,db)
  })

  r.POST("/token", func(c *gin.Context){
    centralsessionserver.TokenRequest(c, db)
  })

  r.GET("/userinfo", func(c *gin.Context){
    centralsessionserver.UserInfoRequest(c, db)
  })

  r.POST("/apps", func(c *gin.Context){
    controlpanel.RegisterApplicationRequest(c, db)
  })

  r.Run()
}