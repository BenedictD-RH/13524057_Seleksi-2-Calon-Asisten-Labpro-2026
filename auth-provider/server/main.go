package main

import (
	"fmt"

  "auth-provider-server/api"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:password@tcp(localhost:5342)/auth-provider-db?charset=utf8mb4&parseTime=True&loc=Local"
var db, db_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func main() {

  if (db_err != nil) {
    fmt.Println(db_err);
  }
  if (db == nil) {
    return;
  } 
  

  r := gin.Default()

  r.POST("/createuser", func(c *gin.Context) {
    api.CreateUserRequest(c, db)
  })

  r.GET("/getusers", func(c *gin.Context) {
    api.GetAllUsersRequest(c, db)
  })

  r.PATCH("/updateuser", func(c *gin.Context) {
    api.UpdateUserRequest(c, db)
  })

  r.DELETE("/deleteuser", func(c *gin.Context) {
    api.DeleteUserRequest(c, db)
  })

    r.POST("/creategroup", func(c *gin.Context) {
    api.CreateGroupRequest(c, db)
  })

  r.GET("/getgroups", func(c *gin.Context) {
    api.GetAllGroupsRequest(c, db)
  })

  r.PATCH("/updategroup", func(c *gin.Context) {
    api.UpdateGroupRequest(c, db)
  })

  r.DELETE("/deletegroup", func(c *gin.Context) {
    api.DeleteGroupRequest(c, db)
  })

  r.Run()
}