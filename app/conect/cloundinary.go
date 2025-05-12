package conect

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary

func ConnectCloudinary() {
	
	KeyName := os.Getenv("KeyName")
	ApiKey := os.Getenv("ApiKey")
	ApiSecret := os.Getenv("ApiSecret")


	cld, err := cloudinary.NewFromParams(KeyName, ApiKey, ApiSecret)
	if err != nil {
		log.Fatal("Fail to connect cloudinary")
	}
	CLD= cld
	fmt.Println("sucessfully to connect")
}
