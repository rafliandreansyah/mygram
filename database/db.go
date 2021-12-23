package database

import (
	"MyGram/model"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "amaterasu"
	dbname   = "mygram_db"
	port     = "5432"

	db  *gorm.DB
	err error
)

func StartDB() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting database =>", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Photo{}, &model.Comment{})
	if err != nil {
		log.Fatal("err:", err.Error())
		return
	}

	fmt.Println("Database connected")

}

func GetDB() *gorm.DB {
	return db
}
