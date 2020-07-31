package API

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"project-backend/API/controllers"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error getting env, not coming through %v" , err)
	}else{
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	fmt.Println("Server run at port :8080")
	server.Run(":8080")
}
