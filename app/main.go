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
	r.GET("/item",handler.GetAllItem)
	r.DELETE("/item/:id",handler.DelItemByID)
	r.Run(":8080")
}