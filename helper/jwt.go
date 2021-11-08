package helper

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

const secretKey = "2isjyr33"

func GenerateToken(id uuid.UUID, username string, email string, age int) (string, error) {
	claim := jwt.MapClaims{
		"id" : id,
		"username": username,
		"email": email,
		"age": strconv.Itoa(age),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(secretKey))
	if  err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(){
	
}
