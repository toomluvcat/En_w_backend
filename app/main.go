package main

import (
	"Render/app/conect"
	"Render/app/handler"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context){
	c.JSON(200,gin.H{"messageg":"Pong"})
}


func main() {
	conect.ConnectDB()
	r:=gin.Default()
	r.POST("/user",handler.CreateUser)
	r.GET("/user/:id",handler.GetUserByID)
	r.PUT("/user/:id",handler.PutUserByID)
	r.POST("/item",handler.CreateItem)
	r.GET("/items/:id",handler.GetAllItem)
	r.GET("/item/:user_id/:item_id",handler.GetAllItem)
	r.DELETE("/item/:id",handler.DelItemByID)

	r.POST("/bookmark",handler.ToggleBookMark)
	r.GET("/tan/:id",handler.GetEventByUserID)
	r.GET("/event",handler.GetAllEvent)
	r.POST("/event",handler.CreateEvent)
	r.PUT("/event/:id",handler.PutEventdByID)
	r.Run(":8080")
}