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
    
    if KeyName == "" || ApiKey == "" || ApiSecret == "" {
        log.Fatal("Cloudinary environment variables are not set properly")
    }

    cld, err := cloudinary.NewFromParams(KeyName, ApiKey, ApiSecret)
    if err != nil {
        log.Fatalf("Fail to connect cloudinary: %v", err)
    }
    CLD = cld
    fmt.Println("Successfully connected to Cloudinary")
}
