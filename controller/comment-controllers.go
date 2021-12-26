package controller

import (
	"MyGram/constant"
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
	"net/http"
	"time"
)

func CreateComment(c *gin.Context) {
	var err error
	var comment model.Comment
	var createCommentReponse struct {
		ID        uuid.UUID `json:"id"`
		Message   string    `json:"message"`
		PhotoID   uuid.UUID `json:"photo_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	contentType := helper.GetContentType(c)
	db := database.GetDB()
	userData := c.MustGet(constant.UserData).(jwt.MapClaims)

	if contentType == constant.JSON {
		err = c.ShouldBindJSON(&comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		err = c.ShouldBind(&comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	userId := uuid.Must(uuid.FromString(userData["id"].(string)))
	comment.UserID = userId

	err = db.Debug().Create(&comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createCommentReponse.ID = comment.ID
	createCommentReponse.Message = comment.Message
	createCommentReponse.UserID = comment.UserID
	createCommentReponse.PhotoID = comment.PhotoID
	createCommentReponse.CreatedAt = comment.CreatedAt

	c.JSON(http.StatusCreated, gin.H{
		"Data": createCommentReponse,
	})

}

func GetSocialMedias(c *gin.Context) {

	type User struct {
		ID       uuid.UUID `json:"id"`
		UserName string    `json:"user_name"`
		Email    string    `json:"email"`
	}

	type SocialMedia struct {
		ID             uuid.UUID `json:"id"`
		Name           string    `json:"name"`
		SocialMediaUrl string    `json:"social_media_url"`
		User           User      `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" valid:"-"`
		UserID         uuid.UUID `json:"user_id" form:"user_id"`
		CreatedAt      time.Time `json:"created_at" form:"created_at"`
	}

	var socialMediasResponse struct {
		SocialMedia []SocialMedia `json:"social_media"`
	}
	var err error
	var socialMedia []SocialMedia

	db := database.GetDB()

	err = db.Debug().Preload(clause.Associations).Find(&socialMedia).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	socialMediasResponse.SocialMedia = socialMedia

	c.JSON(http.StatusOK, gin.H{
		"Data": socialMediasResponse,
	})

}

func DeleteSocialMediaByID(c *gin.Context){

	var err error

	db := database.GetDB()
	socialMediaIdParams := c.Param("socialmedia_id")

	found := db.Debug().First(&model.SocialMedia{ID: uuid.Must(uuid.FromString(socialMediaIdParams))}).RowsAffected
	if found < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "social media not found",
		})
		return
	}

	err = db.Debug().Delete(&model.SocialMedia{ID: uuid.Must(uuid.FromString(socialMediaIdParams))}).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}
