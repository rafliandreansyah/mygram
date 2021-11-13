package controller

import (
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func GetPhotos(c *gin.Context){
	var err error
	var photos []model.Photo
	var user model.User
	getDB := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	//Check content type
	contentType := helper.GetContentType(c)
	if contentType == helper.JSON {
		c.ShouldBindJSON(&photos)
	}else {
		c.ShouldBind(&photos)
	}


	//Check user found
	userId := uuid.Must(uuid.FromString(userData["id"].(string)))
	rowAffected := getDB.First(&user, userId).RowsAffected
	if rowAffected < 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	//Get all photo
	err = getDB.Debug().Find(&photos).Error
	if err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{
			 "message": "error get data",
			 "error": err.Error(),
		 })
		return
	}


	//return all photo
	c.JSON(http.StatusOK, gin.H{
		"photos": photos,
	})

}

func AddPhoto(c *gin.Context){
	var photo model.Photo
	var photoResponse struct{
		ID uuid.UUID
		Title string
		Caption string
		PhotoUrl string
	}
	var user model.User
	var err error
	db := database.GetDB()
	contentType := helper.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	fmt.Println("userdata:", userData)

	if contentType == helper.JSON {
		err = c.ShouldBindJSON(&photo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}else {
		err = c.ShouldBind(&photo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	//Check user found
	userId := uuid.Must(uuid.FromString(userData["id"].(string)))
	rowAffected := db.First(&user, userId).RowsAffected
	if rowAffected < 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	photo.UserID = userId
	fmt.Println("userid:", photo.UserID)

	err = db.Debug().Create(&photo).Error
	if err != nil {
		 c.JSON(http.StatusInternalServerError, gin.H{
			 "error": err.Error(),
		 })
		return
	}

	photoResponse.ID = photo.ID
	photoResponse.Title = photo.Title
	photoResponse.PhotoUrl = photo.PhotoUrl
	photoResponse.Caption = photo.Caption
	c.JSON(http.StatusCreated, gin.H{
		"message": "photo created",
		"photo": photoResponse,
	})

}