package controller

import (
	"MyGram/constant"
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func GetPhotos(c *gin.Context) {
	type Comment struct {
		ID        uuid.UUID `json:"id"`
		Message   string    `json:"message"`
		UserID    uuid.UUID `json:"user_id"`
		PhotoID   uuid.UUID `json:"-"`
		CreatedAt time.Time `json:"created_at"`
	}
	type Photo struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoUrl  string    `json:"photo_url"`
		CreatedAt time.Time `json:"created_at"`
		Comments  []Comment `json:"comments"`
	}

	var err error
	var photos []Photo

	getDB := database.GetDB()

	//Get all photo
	err = getDB.Debug().Preload("Comments").Find(&photos).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//return all photo
	c.JSON(http.StatusOK, gin.H{
		"Data": photos,
	})

}

func AddPhoto(c *gin.Context) {
	var photo model.Photo
	var photoResponse struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoUrl  string    `json:"photo_url"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}
	var user model.User
	var err error
	db := database.GetDB()
	contentType := helper.GetContentType(c)
	userData := c.MustGet(constant.UserData).(jwt.MapClaims)

	if contentType == constant.JSON {
		err = c.ShouldBindJSON(&photo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
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
	photoResponse.UserID = photo.UserID
	photoResponse.CreatedAt = photo.CreatedAt
	c.JSON(http.StatusCreated, gin.H{
		"Data": photoResponse,
	})

}

func GetPhotoByID(c *gin.Context) {
	type User struct {
		ID       uuid.UUID `json:"id"`
		UserName string    `json:"user_name"`
		Email    string    `json:"email"`
	}

	type Comment struct {
		ID        uuid.UUID `json:"id"`
		Message   string    `json:"message"`
		UserID    uuid.UUID `json:"user_id"`
		PhotoID   uuid.UUID `json:"-"`
		CreatedAt time.Time `json:"created_at"`
	}

	type Photo struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoUrl  string    `json:"photo_url"`
		CreatedAt time.Time `json:"created_at"`
		UserID    uuid.UUID `json:"-"`
		User      User      `json:"user"`
		Comments  []Comment `json:"comments"`
	}

	var err error
	var photo Photo

	id := c.Param("photo_id")
	photoID := uuid.Must(uuid.FromString(id))
	db := database.GetDB()

	err = db.Debug().Preload("Comments").Preload("User").Find(&photo, photoID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Data": photo,
	})
}

func EditPhotoByID(c *gin.Context) {
	var err error
	var photo model.Photo
	var photoEdited model.Photo
	var photoResponse struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoUrl  string    `json:"photo_url"`
		UserID    uuid.UUID `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	contentType := helper.GetContentType(c)
	if contentType == constant.JSON {
		err = c.BindJSON(&photoEdited)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		err = c.Bind(&photoEdited)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	id := c.Param("photo_id")
	photoID := uuid.Must(uuid.FromString(id))
	db := database.GetDB()

	data := db.Debug().First(&photo, photoID).RowsAffected
	if data < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "photo not found",
		})
		return
	}

	photo.Title = photoEdited.Title
	photo.Caption = photoEdited.Caption
	photo.PhotoUrl = photoEdited.PhotoUrl
	err = db.Save(&photo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	photoResponse.ID = photo.ID
	photoResponse.Title = photo.Title
	photoResponse.Caption = photo.Caption
	photoResponse.PhotoUrl = photo.PhotoUrl
	photoResponse.UserID = photo.UserID
	photoResponse.UpdatedAt = photo.UpdatedAt

	c.JSON(http.StatusOK, gin.H{
		"Data": photoResponse,
	})

}

func DeletePhotoByID(c *gin.Context) {
	var err error
	var photo model.Photo

	id := c.Param("photo_id")
	photoID := uuid.Must(uuid.FromString(id))

	db := database.GetDB()

	data := db.Debug().First(&photo, photoID).RowsAffected
	if data < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "photo not found",
		})
		return
	}

	err = db.Debug().Delete(&photo, photoID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete photo success",
	})
}
