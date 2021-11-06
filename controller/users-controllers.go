package controller

import (
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var (
	json = "application/json"
	err error
)

func Register(c *gin.Context) {

	var user model.User
	var responseRegister struct {
		ID    uuid.UUID `json:"id"`
		Email string 	`json:"email"`
	}
	var db  = database.GetDB()

	contentType := helper.GetContentType(c)

	//Check content type
	if contentType == json {
		//Bind from json to struct
		err = c.ShouldBindJSON(&user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		err = c.ShouldBind(&user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	//Check age
	if user.Age < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "minimum age must be 8 years old",
		})
		return
	}

	//Create user to database
	err = db.Debug().Create(&user).Error

	// If have error on create user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseRegister.ID = user.ID
	responseRegister.Email = user.Email

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    responseRegister,
	})

}

func Login(c *gin.Context) {

	var user model.User
	var db = database.GetDB()

	contentType := helper.GetContentType(c)
	if contentType == json {
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}else{
		err = c.ShouldBind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	password := user.Password

	// Get check email from database
	err = db.Debug().Where(&model.User{Email: user.Email}).First(&user).Error
	if err != nil {
		fmt.Println("Error get user by email: ",err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "email or password invalid",
		})
		return
	}

	// Check password is match
	isMatch := helper.CompareHashPassword(user.Password, password)
	if !isMatch {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "email or password invalid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Login success",
	})

}
