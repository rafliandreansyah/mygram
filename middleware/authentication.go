package middleware

import (
	"MyGram/constant"
	"MyGram/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		verifyToken, err := helper.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Set(constant.UserData, verifyToken)
		c.Next()
	}
}
