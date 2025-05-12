package conect

import (
	"Render/app/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    HOST := os.Getenv("HOST")
    DATABASE := os.Getenv("DATABASE")
    USERNAME := os.Getenv("USER")
    PASSWORD := os.Getenv("PASSWORD")
    PORT := os.Getenv("PORT")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", HOST, PORT, USERNAME, PASSWORD, DATABASE)

    var err error
    // ใช้ = แทน := เพื่อกำหนดค่าให้กับตัวแปร global
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Fail to connect: %v", err)
    }

    fmt.Println("Connected to the database!")

    err = DB.AutoMigrate(&model.Item{}, &model.User{}, &model.Event{}, &model.Loan{})
    if err != nil {
        log.Fatalf("Fail to migrate: %v", err)
    }

    fmt.Println("Database migrated successfully!")
}