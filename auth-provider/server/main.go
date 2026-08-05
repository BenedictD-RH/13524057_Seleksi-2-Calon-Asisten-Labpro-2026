package main

import (
	"fmt"
	"os"

	"auth-provider-server/api"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db_pass = os.Getenv("AUTH_PROVIDER_DB_PASS")
var dsn = fmt.Sprintf("root:%s@tcp(localhost:5342)/auth-provider-db?charset=utf8mb4&parseTime=True&loc=Local", db_pass)
var db, db_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func main() {

  if (db_err != nil) {
    fmt.Println(db_err);
  }
  if (db == nil) {
    return;
  } 
  

  r := gin.Default()

  r.POST("user/create", func(c *gin.Context) {
    api.CreateUserRequest(c, db)
  })

  r.GET("user/getall", func(c *gin.Context) {
    api.GetAllUsersRequest(c, db)
  })

  r.PATCH("user/update", func(c *gin.Context) {
    api.UpdateUserRequest(c, db)
  })

  r.DELETE("user/delete", func(c *gin.Context) {
    api.DeleteUserRequest(c, db)
  })

    r.POST("group/create", func(c *gin.Context) {
    api.CreateGroupRequest(c, db)
  })

  r.GET("group/getall", func(c *gin.Context) {
    api.GetAllGroupsRequest(c, db)
  })

  r.PATCH("group/update", func(c *gin.Context) {
    api.UpdateGroupRequest(c, db)
  })

  r.DELETE("group/delete", func(c *gin.Context) {
    api.DeleteGroupRequest(c, db)
  })

  r.POST("group/adduser", func(c *gin.Context) {
    api.AddUserToGroup(c, db)
  })

  r.DELETE("group/removeuser", func(c *gin.Context) {
    api.RemoveUserFromGroup(c, db)
  })

  r.Run()
}