package centralsessionserver

import (
	"auth-provider-server/api/models"
	"auth-provider-server/api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessTokenPayload struct {
	Token string `json:"access_token" binding:"required"`
}

func FindAccessToken(access_tokens *[]models.AccessToken, access_token string) *models.AccessToken {
	for i := 0; i < len(*access_tokens); i++ {
		if utility.CheckPasswordHash(access_token, (*access_tokens)[i].TokenHash) {
			return &(*access_tokens)[i]
		}
	}
	return nil
}

func UserInfoRequest(c *gin.Context, db *gorm.DB) {
	var tokenPayload AccessTokenPayload

	if err := c.ShouldBindJSON(&tokenPayload); err != nil {
		UnauthorizedResponse(c)
		return
	}

	var access_tokens []models.AccessToken
	db.Where("status = ?", "Active").Find(&access_tokens)
	accessTokenModel := FindAccessToken(&access_tokens, tokenPayload.Token)
	if accessTokenModel == nil {
		UnauthorizedResponse(c)
		return
	}
	accessTokenModel.UpdateStatus(db)
	if !accessTokenModel.IsValid() {
		UnauthorizedResponse(c)
		return
	}

	var user models.User
	var users []models.User
	db.Where("id = ?", accessTokenModel.UserId).Find(&users)
	if len(users) <= 0 {
		UnauthorizedResponse(c)
		return
	}
	user = users[0]

	c.JSON(http.StatusOK, gin.H{
		"central_session_id": accessTokenModel.SsoSessionId,
		"user_id":            user.ID,
		"name":               user.Name,
		"email":              user.Email,
		"groups":             *(user.GetGroupsUUID(db)),
	})
}
