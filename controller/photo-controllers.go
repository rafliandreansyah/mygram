package controller

import (
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPhotos(c *gin.Context){
	var err error
	var photos []model.Photo
	getDB := database.GetDB()

	contentType := helper.GetContentType(c)
	if contentType == helper.JSON {
		c.ShouldBindJSON(&photos)
	}else {
		c.ShouldBind(&photos)
	}

	err = getDB.Debug().Find(&photos).Error
	if err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{
			 "message": "error get data",
			 "error": err.Error(),
		 })
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photos": photos,
	})

}