package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAppToAuthProviderRequest(c *gin.Context, db *gorm.DB) {
	auth_server_url := os.Getenv("AUTH_SERVER_URL")
	launch_uri := os.Getenv("FRONTEND_URI")
	logout_notification_url := os.Getenv("BACKEND_URI") + "/internal/logout"
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	app_name := os.Getenv("APP_NAME")

	if auth_server_url == "" || launch_uri == "" || client_id == "" || client_secret == "" || app_name == "" {
		fmt.Println("Missing environment variables")
		//InternalServerError
		return
	}

	payload := map[string]string{"name": app_name,
		"client_id":     client_id,
		"client_secret": client_secret,
		"launch_url":  launch_uri,
		"logout_notification_url": logout_notification_url}
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/apps", auth_server_url), bytes.NewReader(jsonData))

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		InternalServerErrorResponse(c)
		return
	}

	c.JSON(resp.StatusCode, resp.Body)
}
