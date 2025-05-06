package conect

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Errorf("fail to load env: %v",err)
	}
	HOST := os.Getenv("HOST")
	DATABASE :=os.Getenv("DATABASE")
	USERNAME := os.Getenv("USER")
	PASSWORD := os.Getenv("PASSWORD")
	PORT  := os.Getenv("PORT")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", HOST, PORT, USERNAME, PASSWORD, DATABASE)
	DB,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		fmt.Errorf("Fail to connect: %v",err)
	}	

}
