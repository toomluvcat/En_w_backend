package conect

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var CLD *cloudinary.Cloudinary

func ConnectCloudinary() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Fail to load env")
	}

	KeyName := os.Getenv("KeyName")
	ApiKey := os.Getenv("ApiKey")
	ApiSecret := os.Getenv("ApiSecret")


	CLD, err = cloudinary.NewFromParams(KeyName, ApiKey, ApiSecret)
	if err != nil {
		log.Fatal("Fail to connect cloudinary")
	}
}
