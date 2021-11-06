package helper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	pass := []byte(p)
	cost := 11

	passHash, err := bcrypt.GenerateFromPassword(pass, cost)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func CompareHashPassword(hashPassword string, passwordFromUserLogin string) bool{
	hashPass := []byte(hashPassword)
	passCompare := []byte(passwordFromUserLogin)
	err := bcrypt.CompareHashAndPassword(hashPass, passCompare)
	if err != nil {
		fmt.Println("Error compare password:", err.Error())
		return false
	}
	return true
}