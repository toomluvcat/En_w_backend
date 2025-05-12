package main

import (
	"Render/app/conect"
	"Render/app/handler"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context){
	c.JSON(200,gin.H{"messageg":"Pong"})
}


func main() {
	conect.ConnectCloudinary()
	conect.ConnectDB()
	r:=gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
	r.POST("/user",handler.CreateUser)
	r.GET("/user/:id",handler.GetUserByID)
	r.PUT("/user/:id",handler.PutUserByID)
	r.POST("/item",handler.CreateItem)
	r.GET("/items/:id",handler.GetAllItem)
	r.GET("/item/:user_id/:item_id",handler.GetAllItem)
	r.DELETE("/item/:id",handler.DelItemByID)
	r.GET("/admin/item",handler.GetAllItemByAdmin)
	r.GET("/admin/item/:id",handler.GetItemByIDAdmin)
	r.POST("/bookmark",handler.ToggleBookMark)
	r.GET("/event/:id",handler.GetEventByUserID)
	r.GET("/event",handler.GetAllEvent)
	r.POST("/event",handler.CreateEvent)
	r.PUT("/admin/item/img/:id",handler.PutItemByID)
	r.PUT("/admin/item/:id",handler.PutItemByIDNoImage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback เผื่อใช้ local
	}
	r.Run(":"+port)
}