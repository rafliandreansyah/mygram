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

func PostSocialMedia(c *gin.Context) {

	var socialMediaResponse struct {
		ID             uuid.UUID `json:"id"`
		Name           string    `json:"name"`
		SocialMediaUrl string    `json:"social_media_url"`
		UserID         uuid.UUID `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
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
	} else {
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

func EditSocialMediaByID(c *gin.Context) {
	var err error
	var socialMedia model.SocialMedia
	var socialMediaEdit model.SocialMedia
	var socialMediaResponse struct {
		ID             uuid.UUID `json:"id"`
		Name           string    `json:"name"`
		SocialMediaUrl string    `json:"social_media_url"`
		UserID         uuid.UUID `json:"user_id"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	socialMediaID := c.Param("socialmedia_id")
	contentType := helper.GetContentType(c)

	db := database.GetDB()

	//Check content type
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

	//Check social media found and get data
	sosmedConvertID, err := uuid.FromString(socialMediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	socialMediaEdit.ID = uuid.Must(sosmedConvertID, err)
	socialMediaData := db.Debug().First(&socialMediaEdit).RowsAffected
	if socialMediaData < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "social media not found",
		})
		return
	}

	//Edit social media
	socialMediaEdit.Name = socialMedia.Name
	socialMediaEdit.SocialMediaUrl = socialMedia.SocialMediaUrl
	err = db.Save(&socialMediaEdit).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	socialMediaResponse.ID = socialMediaEdit.ID
	socialMediaResponse.Name = socialMediaEdit.Name
	socialMediaResponse.SocialMediaUrl = socialMediaEdit.SocialMediaUrl
	socialMediaResponse.UserID = socialMediaEdit.UserID
	socialMediaResponse.UpdatedAt = socialMediaEdit.UpdatedAt

	c.JSON(http.StatusOK, gin.H{
		"Data": socialMediaResponse,
	})

}

func DeleteSocialMediaByID(c *gin.Context) {

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
