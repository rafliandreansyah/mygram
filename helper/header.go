package helper

import "github.com/gin-gonic/gin"

func GetContentType(c *gin.Context) string {
	contentType := c.Request.Header.Get("Content-Type")
	return contentType
}
