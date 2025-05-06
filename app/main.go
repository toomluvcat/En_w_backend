package main

import (
	"render/app/conect"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context){
	c.JSON(200,gin.H{"messageg":"Pong"})
}


func main() {
	conect.ConnectDB()
	r:=gin.Default()
	
	r.GET("/",Ping)

	r.Run()
}