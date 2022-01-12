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

func GetComment(c *gin.Context) {
	type Photo struct {
		ID       uuid.UUID `json:"id"`
		Title    string    `json:"title"`
		Caption  string    `json:"caption"`
		PhotoUrl string    `json:"photo_url"`
		UserID   uuid.UUID `json:"user_id"`
	}
	type Comment struct {
		ID        uuid.UUID   `json:"id"`
		Message   string      `json:"message"`
		PhotoID   uuid.UUID   `json:"photo_id"`
		UpdatedAt time.Time   `json:"updated_at"`
		CreatedAt time.Time   `json:"created_at"`
		Photo     Photo `json:"photo"`
	}
	var err error
	var comments []Comment

	db := database.GetDB()

	err = db.Debug().Preload("Photo").Find(&comments).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Data": comments,
	})

}

func CreateComment(c *gin.Context) {
	var err error
	var comment model.Comment
	var createCommentResponse struct {
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

	createCommentResponse.ID = comment.ID
	createCommentResponse.Message = comment.Message
	createCommentResponse.UserID = comment.UserID
	createCommentResponse.PhotoID = comment.PhotoID
	createCommentResponse.CreatedAt = comment.CreatedAt

	c.JSON(http.StatusCreated, gin.H{
		"Data": createCommentResponse,
	})

}

func EditCommentByID(c *gin.Context) {
	var err error
	var comment model.Comment
	var commentEdit model.Comment
	var commentResponse struct {
		ID        uuid.UUID `json:"id"`
		Message   string    `json:"message"`
		PhotoID   uuid.UUID `json:"photo_id"`
		UserID    uuid.UUID `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	db := database.GetDB()

	commentID := c.Param("comment_id")
	contentType := helper.GetContentType(c)

	//Check content type
	if contentType == constant.JSON {
		err = c.BindJSON(&comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
	} else {
		err = c.Bind(&comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
	}

	//Check comment found and get data
	commentEdit.ID = uuid.Must(uuid.FromString(commentID))
	commentData := db.Debug().First(&commentEdit).RowsAffected

	if commentData < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "comment not found",
		})
		return
	}

	//Edit comment
	commentEdit.Message = comment.Message
	err = db.Save(&commentEdit).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	commentResponse.ID = commentEdit.ID
	commentResponse.Message = commentEdit.Message
	commentResponse.PhotoID = commentEdit.PhotoID
	commentResponse.UserID = commentEdit.UserID
	commentResponse.UpdatedAt = commentEdit.UpdatedAt
	c.JSON(http.StatusOK, gin.H{
		"Data": commentResponse,
	})

}

func DeleteCommentByID(c *gin.Context) {
	var err error

	id := c.Param("comment_id")
	db := database.GetDB()

	// Check comment
	comment := db.Debug().First(&model.Comment{
		ID: uuid.Must(uuid.FromString(id)),
	}).RowsAffected
	if comment < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "comment not found",
		})
		return
	}

	//Delete comment
	err = db.Debug().Delete(&model.Comment{
		ID: uuid.Must(uuid.FromString(id)),
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "comment deleted",
	})

}
