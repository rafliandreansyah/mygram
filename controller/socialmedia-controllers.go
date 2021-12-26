package controller

import (
	"MyGram/constant"
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func PostSocialMedia(c *gin.Context){

	var socialMediaResponse struct {
		ID uuid.UUID `json:"id"`
		Name string `json:"name"`
		SocialMediaUrl string `json:"social_media_url"`
		UserID uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	var socialMedia model.SocialMedia
	var err error

	contentType := helper.GetContentType(c)
	userData := c.MustGet(constant.UserData).(jwt.MapClaims)
	db := database.GetDB()

	if contentType == constant.JSON {
		err = c.BindJSON(&socialMedia)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
	}else {
		err = c.Bind(&socialMedia)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
	}

	userID := uuid.Must(uuid.FromString(userData["id"].(string)))
	socialMedia.UserID = userID

	err = db.Debug().Create(&socialMedia).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	socialMediaResponse.ID = socialMedia.ID
	socialMediaResponse.SocialMediaUrl = socialMedia.SocialMediaUrl
	socialMediaResponse.Name = socialMedia.Name
	socialMediaResponse.CreatedAt = socialMedia.CreatedAt
	socialMediaResponse.UserID = socialMedia.UserID

	c.JSON(http.StatusOK, gin.H{
		"Data": socialMediaResponse,
	})

}