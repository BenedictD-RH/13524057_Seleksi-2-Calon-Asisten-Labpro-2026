package auth

import (
	"app-a-backend/api/models"
	"app-a-backend/api/utility"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var code_verifier_exp_duration = 2 * time.Minute

func CreateCodeVerifier() (models.CodeVerifier, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return models.CodeVerifier{}, err
	}
	state, err := utility.CryptoRandString(32)
	if err != nil {
		return models.CodeVerifier{}, err
	}
	code_verifier, err := utility.CryptoRandString(60)
	if err != nil {
		return models.CodeVerifier{}, err
	}
	return models.CodeVerifier{
		ID:           datatypes.BinUUID(newUUID),
		State:        state,
		CodeVerifier: code_verifier,
		Status:       "Active",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(code_verifier_exp_duration),
	}, nil
}

func LoginRequest(c *gin.Context, db *gorm.DB) {
	auth_server_url := os.Getenv("AUTH_SERVER_URL")
	redirect_uri := os.Getenv("REDIRECT_URI")
	client_id := os.Getenv("CLIENT_ID")

	if auth_server_url == "" || redirect_uri == "" || client_id == "" {
		fmt.Println("Missing environment variables")
		InternalServerErrorResponse(c)
		return
	}

	codeVerifierModel, err := CreateCodeVerifier()
	if err != nil {
		InternalServerErrorResponse(c)
		return
	}

	if err := db.Create(&codeVerifierModel).Error; err != nil {
		InternalServerErrorResponse(c)
		return
	}

	code_challenge, err := utility.HashString(codeVerifierModel.CodeVerifier)
	fmt.Println(code_challenge)
	if err != nil {
		InternalServerErrorResponse(c)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/authorize?redirect_uri=%s&client_id=%s&state=%s&code_challenge=%s",
		auth_server_url, redirect_uri, client_id, codeVerifierModel.State, code_challenge))
}
