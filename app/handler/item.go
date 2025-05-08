package handler

import (
	"Render/app/conect"
	"Render/app/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context){
	type Req struct{
		Name string
		Quantity int
	}

	var req Req
	if err := c.BindJSON(&req);err!=nil{
		c.Status(400)
		return
	}

	item :=model.Item{
		Name:req.Name,
		MaxQuantity: req.Quantity,
		CurrentQuantity: req.Quantity,
	}

	 result:= conect.DB.Create(&item);
	 if result.Error!=nil{
		c.Status(500)
		return	
	}
	c.JSON(200,result)

}

func GetAllItem(c *gin.Context){
	var items []model.Item

	if err:= conect.DB.Find(&items).Error;err!=nil{
		c.Status(500)
	}
	c.JSON(200,items)
	
}


func DelItemByID(c *gin.Context){
	id := c.Param("id")

	if err:= conect.DB.Delete(&model.Item{},id).Error;err!=nil{
		c.Status(500)
	}
	c.Status(201)
}


func PutItemByID(c *gin.Context) {
	id := c.Param("id")

	var req model.Item
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	if err := conect.DB.Model(&model.Item{}).Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to update item: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "Item updated successfully"})
}
