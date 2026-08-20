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

	mq_pass := os.Getenv("MESSAGE_QUEUE_DB_PASS")
	dsn_mq := fmt.Sprintf("root:%s@tcp(localhost:5343)/queue?charset=utf8mb4&parseTime=True&loc=Local", mq_pass)
	mq, db_err := gorm.Open(mysql.Open(dsn_mq), &gorm.Config{})
	if db_err != nil {
		fmt.Println(db_err)
		return
	}
	if mq == nil {
		return
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  	  []string{os.Getenv("AUTH_PORTAL_URI"), 
								   os.Getenv("ADMIN_CONSOLE_URI"), 
								   "https://localhost:8692", "https://localhost:8691"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/users", func(c *gin.Context) {
		controlpanel.CreateUserRequest(c, db)
	})

	r.GET("/users", func(c *gin.Context) {
		controlpanel.GetAllUsersRequest(c, db)
	})

	r.GET("/users/fields", controlpanel.GetUserFields)

	r.PATCH("/users", func(c *gin.Context) {
		controlpanel.UpdateUserRequest(c, db, mq)
	})

	r.DELETE("/users", func(c *gin.Context) {
		controlpanel.DeleteUserRequest(c, db, mq)
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
		controlpanel.DeleteGroupRequest(c, db, mq)
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
		controlpanel.RemoveUserFromGroup(c, db, mq)
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
		controlpanel.RemoveUserFromGroup(c, db, mq)
	})

	r.POST("/groups/apps", func(c *gin.Context) {
		controlpanel.AddAppGroupPolicyRequest(c, db)
	})

	r.POST("/groups/apps/query", func(c *gin.Context) {
		controlpanel.GetAppPoliciesFromGroupRequest(c, db)
	})

	r.POST("/groups/apps/query/complement", func(c *gin.Context) {
		controlpanel.GetAppsNotInGroupPoliciesRequest(c, db)
	})

	r.GET("/groups/apps/fields", controlpanel.GetAppPoliciesFieldFromGroupRequest)

	r.PATCH("/groups/apps", func(c *gin.Context) {
		controlpanel.UpdateAppGroupPolicyRequest(c, db, mq)
	})

	r.DELETE("/groups/apps", func(c *gin.Context) {
		controlpanel.RemoveAppGroupPolicyRequest(c, db, mq)
	})

	r.POST("/login", func(c *gin.Context) {
		centralsessionserver.LoginRequest(c, db, mq)
	})

	r.GET("/authorize", func(c *gin.Context) {
		centralsessionserver.AuthorizeRequest(c, db)
	})

	r.GET("/authorize/administrator", func(c *gin.Context) {
		centralsessionserver.AuthorizeAdminRequest(c, db)
	})

	r.POST("/token", func(c *gin.Context) {
		centralsessionserver.TokenRequest(c, db)
	})

	r.GET("/userinfo", func(c *gin.Context) {
		centralsessionserver.UserInfoRequest(c, db)
	})

	r.GET("/userinfo/central", func(c *gin.Context) {
		centralsessionserver.UserInfoInternalRequest(c, db)
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

	r.POST("/apps/groups", func(c *gin.Context) {
		controlpanel.AddAppGroupPolicyRequest(c, db)
	})

	r.POST("/apps/groups/query", func(c *gin.Context) {
		controlpanel.GetGroupPoliciesFromAppRequest(c, db)
	})

	r.POST("/apps/groups/query/complement", func(c *gin.Context) {
		controlpanel.GetGroupsNotInAppPoliciesRequest(c, db)
	})

	r.GET("/apps/groups/fields", controlpanel.GetGroupPoliciesFieldFromAppRequest)

	r.PATCH("/apps/groups", func(c *gin.Context) {
		controlpanel.UpdateAppGroupPolicyRequest(c, db, mq)
	})

	r.DELETE("/apps/groups", func(c *gin.Context) {
		controlpanel.RemoveAppGroupPolicyRequest(c, db, mq)
	})

	r.POST("/apps/uri", func(c *gin.Context) {
		controlpanel.RegisterAppURIRequest(c, db)
	})

	r.POST("/apps/uri/query", func(c *gin.Context) {
		controlpanel.GetAllAppURIRequest(c, db)
	})

	r.GET("/apps/uri/fields", controlpanel.GetURIFields)

	r.DELETE("/apps/uri", func(c *gin.Context) {
		controlpanel.RemoveAppURIRequest(c, db)
	})

	r.POST("/logout", func(c *gin.Context) {
		centralsessionserver.LogoutRequest(c, db, mq)
	})

	r.GET("/session", func(c *gin.Context) {
		centralsessionserver.GetSessionInfoRequest(c, db)
	})

	r.GET("/session/use", func(c *gin.Context) {
		centralsessionserver.PickSessionRequest(c, db)
	})

	r.GET("/return", func(c *gin.Context) {
		centralsessionserver.ReturnToClientPage(c, db)
	})

	r.RunTLS(":8080", "./certs/localhost.pem", "./certs/localhost-key.pem")
}
