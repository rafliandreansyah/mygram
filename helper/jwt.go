package helper

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
)

const secretKey = "2isjyr33"

func GenerateToken(id uuid.UUID, username string, email string, age int) (string, error) {
	claim := jwt.MapClaims{
		"id" : id,
		"username": username,
		"email": email,
		"age": strconv.Itoa(age),
		"time": time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(secretKey))
	fmt.Println("Secret key:", []byte(secretKey))
	if  err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(c *gin.Context)(interface{}, error){
	errAuthorization := errors.New("not authenticated")
	headerToken := c.Request.Header.Get("Authorization")
	hasBearer := strings.HasPrefix(headerToken, "Bearer")

	if !hasBearer {
		return nil, errAuthorization
	}

	tokenString := strings.Split(headerToken, " ")[1]
	fmt.Println("Token:", tokenString)

	token, err:= jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errAuthorization
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errAuthorization
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}else {
		return nil, errAuthorization
	}


}
