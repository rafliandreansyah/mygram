package controller

import (
	"MyGram/constant"
	"MyGram/database"
	"MyGram/helper"
	"MyGram/model"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func Register(c *gin.Context) {
	var err error
	var user model.User
	var responseRegister struct {
		ID        uuid.UUID `json:"id"`
		Age       int       `json:"age"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
	}
	var db = database.GetDB()

	contentType := helper.GetContentType(c)

	//Check content type
	if contentType == constant.JSON {
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

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
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
	responseRegister.Age = user.Age
	responseRegister.Password = user.Password
	responseRegister.Username = user.UserName
	responseRegister.CreatedAt = user.CreatedAt

	c.JSON(http.StatusCreated, gin.H{
		"Data":    responseRegister,
	})

}

func Login(c *gin.Context) {
	var err error
	var user model.User
	var db = database.GetDB()

	contentType := helper.GetContentType(c)
	if contentType == constant.JSON {
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
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
		fmt.Println("Error get user by email: ", err.Error())
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

	token, err := helper.GenerateToken(user.ID, user.UserName, user.Email, user.Age)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func DeleteUser(c *gin.Context) {

	var user model.User
	db := database.GetDB()
	userData := c.MustGet(constant.UserData).(jwt.MapClaims)

	//Check user found
	userId := uuid.Must(uuid.FromString(userData["id"].(string)))
	rowAffected := db.First(&user, userId).RowsAffected
	if rowAffected < 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	err := db.Select("Photos").Delete(&model.User{ID: userId}).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed delete account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})

}
