package handler

import (
	"Render/app/conect"
	"Render/app/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context) {
	type Req struct {
		Name     string
		Category string
		Quantity int
	}

	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.Status(400)
		return
	}

	item := model.Item{
		Name:            req.Name,
		Category:        req.Category,
		MaxQuantity:     req.Quantity,
		CurrentQuantity: req.Quantity,
	}

	result := conect.DB.Create(&item)
	if result.Error != nil {
		c.Status(500)
		return
	}
	c.JSON(200, result)

}

func ToggleBookMark(c *gin.Context) {
	type Req struct {
		UserID uint
		ItemID uint
	}

	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	var user model.User
	if err := conect.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var item model.Item
	if err := conect.DB.First(&item, req.ItemID).Error; err != nil {
		c.JSON(404, gin.H{"error": "item not found"})
		return
	}

	var bookmarked bool
	for _, b := range user.BookmarksItems {
		if b.ID == req.ItemID {
			bookmarked = true
			break
		}
	}
	if bookmarked {
		if err := conect.DB.Model(&user).Association("BookmarksItems").Delete(&item); err != nil {
			c.JSON(500, gin.H{"error": "Fail to delete bookmark"})
			return
		}
		c.Status(200)
	} else {
		if err := conect.DB.Model(&user).Association("BookmarksItems").Append(&item); err != nil {
			c.JSON(500, gin.H{"error": "Fail to add bookmark"})
			return
		}
		c.Status(201)
	}
}

func GetItemByID(c *gin.Context) {
	itemID := c.Param("item_id")
	userID := c.Param("user_id")

	var item model.Item
	if err := conect.DB.Where("id=?", itemID).First(&item).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to load itemID"})
		return
	}

	var count int64
	result := conect.DB.Table("bookmarks").Where("item_id = ? AND user_id = ?", itemID, userID).Count(&count)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Fail to fetch bookmarks"})
		return
	}

	type Response struct {
		ItemID	uint
		Name   string
		Description string
		CurrentQuantity int
		MaxQuantity int
		Bookmarked bool
	}
	isBookmark := false
	if count>0 {
		isBookmark = true
	}
	res :=Response{
		ItemID: item.ID,
		Name: item.Name,
		Description: item.Description,
		CurrentQuantity: item.CurrentQuantity,
		MaxQuantity: item.MaxQuantity,
		Bookmarked: isBookmark,
	}
	c.JSON(200,res)

}

func GetAllItem(c *gin.Context) {

	userID := c.Param("user_id")

	type ItemResponse struct {
		ItemID          uint
		Name            string
		Description     string
		Bookmarks       bool
		Category        string
		MaxQuantity     int
		CurrentQuantity int
	}

	var items []model.Item

	if err := conect.DB.Find(&items).Error; err != nil {
		c.Status(500)
		return
	}
	var bookmarksIDs []uint
	if err := conect.DB.Table("bookmarks").Where("user_id = ?", userID).Pluck("item_id", &bookmarksIDs).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to Fetch bookmarks"})
		return
	}

	bookmarksMap := make(map[uint]struct{}, len(bookmarksIDs))
	for _, ID := range bookmarksIDs {
		bookmarksMap[ID] = struct{}{}
	}

	var res []ItemResponse
	for _, item := range items {
		_, isBookmark := bookmarksMap[item.ID]
		res = append(res, ItemResponse{
			ItemID:          item.ID,
			Name:            item.Name,
			Category:        item.Category,
			Description:     item.Description,
			Bookmarks:       isBookmark,
			MaxQuantity:     item.MaxQuantity,
			CurrentQuantity: item.CurrentQuantity,
		})
	}
	c.JSON(200, res)

}

func DelItemByID(c *gin.Context) {
	id := c.Param("id")

	if err := conect.DB.Delete(&model.Item{}, id).Error; err != nil {
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
