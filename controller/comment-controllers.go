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
