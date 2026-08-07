package auth

import (
	"app-a-backend/api/utility"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func LoginRequest(c *gin.Context, db *gorm.DB) {
	auth_server_url := os.Getenv("AUTH_SERVER_URL")
	redirect_uri := os.Getenv("REDIRECT_URI")
	client_id := os.Getenv("CLIENT_ID")
	state, err := utility.CryptoRandString(32)
	if (err != nil) {
		//InternalServerError
		return
	}
	code_verifier, err := utility.CryptoRandString(64)
	if (err != nil) {
		//InternalServerError
		return
	}
	code_challenge, err := utility.HashString(code_verifier)
	if (err != nil) {
		//InternalServerError
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/authorize?redirect_uri=%s&client_id=%s&state=%s&code_challenge=%s", 
			   auth_server_url, redirect_uri, client_id, state, code_challenge))
}