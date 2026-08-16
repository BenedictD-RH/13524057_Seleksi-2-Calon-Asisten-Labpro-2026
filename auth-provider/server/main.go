package main

import (
	"fmt"
	"os"

	centralsessionserver "auth-provider-server/api/central-session-server"
	controlpanel "auth-provider-server/api/control-panel"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Remove for docker
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Missing .env file: " + err.Error())
		return
	}
	// Remove for docker
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

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{os.Getenv("AUTH_PORTAL_URI"),
			os.Getenv("ADMIN_CONSOLE_URI"),
			"http://localhost:8692",
			"http://localhost:8961",
		},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	r.POST("/users", func(c *gin.Context) {
		controlpanel.CreateUserRequest(c, db)
	})

	r.GET("/users", func(c *gin.Context) {
		controlpanel.GetAllUsersRequest(c, db)
	})

	r.GET("/users/fields", controlpanel.GetUserFields)

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

	r.GET("/groups/fields", controlpanel.GetGroupFields)

	r.PATCH("/groups", func(c *gin.Context) {
		controlpanel.UpdateGroupRequest(c, db)
	})

	r.DELETE("/groups", func(c *gin.Context) {
		controlpanel.DeleteGroupRequest(c, db)
	})

	r.POST("/users/groups", func(c *gin.Context) {
		controlpanel.AddUserToGroup(c, db)
	})

	r.POST("/users/groups/query", func(c *gin.Context) {
		controlpanel.GetGroupsFromUserRequest(c, db)
	})
  
  r.POST("/users/groups/query/complement", func(c *gin.Context) {
		controlpanel.GetGroupsUserIsNotInRequest(c, db)
	})

  r.GET("/users/groups/fields", controlpanel.GetGroupsFromUserFieldsRequest)

	r.DELETE("/users/groups", func(c *gin.Context) {
		controlpanel.RemoveUserFromGroup(c, db)
	})

	r.POST("/groups/users", func(c *gin.Context) {
		controlpanel.AddUserToGroup(c, db)
	})

  r.POST("/groups/users/query", func(c *gin.Context) {
		controlpanel.GetUsersFromGroupRequest(c, db)
	})

  r.POST("/groups/users/query/complement", func(c *gin.Context) {
		controlpanel.GetUsersNotInGroupRequest(c, db)
	})

  r.GET("/groups/users/fields", controlpanel.GetUsersFromGroupFieldsRequest)

	r.DELETE("/groups/users", func(c *gin.Context) {
		controlpanel.RemoveUserFromGroup(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		centralsessionserver.LoginRequest(c, db)
	})

	r.GET("/authorize", func(c *gin.Context) {
		centralsessionserver.AuthorizeRequest(c, db)
	})

	r.POST("/token", func(c *gin.Context) {
		centralsessionserver.TokenRequest(c, db)
	})

	r.GET("/userinfo", func(c *gin.Context) {
		centralsessionserver.UserInfoRequest(c, db)
	})

	r.POST("/apps", func(c *gin.Context) {
		controlpanel.RegisterApplicationRequest(c, db)
	})

	r.GET("/apps", func(c *gin.Context) {
		controlpanel.GetAllApplicationsRequest(c, db)
	})

	r.GET("/apps/fields", controlpanel.GetApplicationFields)

	r.PATCH("/apps", func(c *gin.Context) {
		controlpanel.UpdateAppRequest(c, db)
	})

	r.DELETE("/apps", func(c *gin.Context) {
		controlpanel.RemoveAppRequest(c, db)
	})

	r.Run()
}
