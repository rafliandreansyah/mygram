package main

import (
	"MyGram/database"
	"MyGram/router"
	"fmt"
)

func main() {

	database.StartDB()

	r := router.MainRouter()

	err := r.Run(":8081")
	if err != nil {
		fmt.Println("Error run server")
		return
	}

}
