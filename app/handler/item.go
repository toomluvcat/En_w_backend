package handler

import (
	"Render/app/conect"
	"Render/app/model"
	"fmt"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "Fail to load file"})
		return
	}
	defer file.Close()

	url, err := UploadToCld(c, file, fileHeader.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload to Cloudinary"})
		return
	}

	itemDataJson := c.PostForm("itemData")

	type ItemData struct {
		Name            string
		Description     string
		Category        string
		CurrentQuantity int
		MaxQuantity     int
	}

	var itemData ItemData
	if err := json.Unmarshal([]byte(itemDataJson), &itemData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	item := model.Item{
		Name:            itemData.Name,
		ImageUrl:        url,
		Description:     itemData.Description,
		Category:        itemData.Category,
		MaxQuantity:     itemData.MaxQuantity,
		CurrentQuantity: itemData.CurrentQuantity,
	}

	result := conect.DB.Create(&item)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database insert failed"})
		return
	}
	c.JSON(200, gin.H{"url": url})
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
		ItemID          uint
		Name            string
		Description     string
		CurrentQuantity int
		MaxQuantity     int
		ImageUrl        string
		Bookmarked      bool
	}
	isBookmark := false
	if count > 0 {
		isBookmark = true
	}
	res := Response{
		ItemID:          item.ID,
		Name:            item.Name,
		Description:     item.Description,
		ImageUrl:        item.ImageUrl,
		CurrentQuantity: item.CurrentQuantity,
		MaxQuantity:     item.MaxQuantity,
		Bookmarked:      isBookmark,
	}
	c.JSON(200, res)

}

func GetAllItem(c *gin.Context) {

	userID := c.Param("user_id")

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	limit := 30
	offset := (page - 1) * limit

	type ItemResponse struct {
		ItemID          uint
		Name            string
		Description     string
		ImageUrl        string
		Bookmarks       bool
		Category        string
		MaxQuantity     int
		CurrentQuantity int
	}

	var items []model.Item

	if err := conect.DB.Find(&items).Offset(offset).Limit(limit).Error; err != nil {
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
			ImageUrl:        item.ImageUrl,
			Bookmarks:       isBookmark,
			MaxQuantity:     item.MaxQuantity,
			CurrentQuantity: item.CurrentQuantity,
		})
	}
	c.JSON(200, res)

}

func GetAllItemByAdmin(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	limit := 30
	offset := (page - 1) * limit

	type ItemResponse struct {
		ItemID          uint
		Name            string
		Description     string
		ImageUrl        string
		Bookmarks       bool
		Category        string
		MaxQuantity     int
		CurrentQuantity int
	}

	var items []model.Item

	if err := conect.DB.Find(&items).Offset(offset).Limit(limit).Error; err != nil {
		c.Status(500)
		return
	}

	var res []ItemResponse
	for _, item := range items {
		res = append(res, ItemResponse{
			ItemID:          item.ID,
			Name:            item.Name,
			Category:        item.Category,
			Description:     item.Description,
			ImageUrl:        item.ImageUrl,
			MaxQuantity:     item.MaxQuantity,
			CurrentQuantity: item.CurrentQuantity,
		})
	}
	c.JSON(200, res)

}

func GetItemByIDAdmin(c *gin.Context) {
	id := c.Param("id")

	type LoanItem struct {
		ItemID   uint
		Name     string
		Quantity int
	}

	type Response struct {
		EventID    uint
		UserName   string
		UserID     uint
		CreatedAt  time.Time
		ApprovedAt time.Time
		Status     string
		Loan       []LoanItem
	}

	var event []model.Event
	if err := conect.DB.Preload("User").
		Joins("JOIN loans ON loans.event_id = events.id").Where("loans.item.id=?", id).Group("events.id").Error; err != nil {
		c.JSON(500, gin.H{"message": "Fail to fetch event"})
		return
	}
	var res []Response
	for _, e := range event {
		var loanItems []LoanItem
		for _, l := range e.Loans {
			loanItems = append(loanItems, LoanItem{
				ItemID:   l.ItemID,
				Name:     l.Item.Name,
				Quantity: l.Quantity,
			})
		}
		var ApprovedAt time.Time
		if e.Status != "Pending" {
			ApprovedAt = e.CreatedAt
		}
		res = append(res, Response{
			EventID:    e.ID,
			UserName:   e.User.Name,
			UserID:     e.UserID,
			CreatedAt:  e.CreatedAt,
			ApprovedAt: ApprovedAt,
			Status:     e.Status,
			Loan:       loanItems,
		})
	}

	type ItemResponse struct {
		ID              uint
		Name            string
		Description     string
		Category        string
		CurrentQuantity int
		MaxQuantity     int
		ImageUrl        string
	}

	var itemResponse ItemResponse
	if err := conect.DB.Model(&model.Item{}).First(&itemResponse, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to fetch itemdata"})
		return
	}

	c.JSON(200, gin.H{"item": itemResponse, "event": event})
}
func DelItemByID(c *gin.Context) {
	id := c.Param("id")

	var item model.Item
	if err := conect.DB.First(&item, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Item not found: " + err.Error()})
		return
	}

	if item.CurrentQuantity != item.MaxQuantity {
		c.JSON(400, gin.H{"error": "Cannot delete item: some are still borrowed"})
		return
	}

	if err := conect.DB.Delete(&item).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete item: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Item deleted successfully"})
}

func PutItemByID(c *gin.Context) {
	id := c.Param("id")

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil && file == nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}

	itemDataJSON := c.PostForm("itemData")

	var itemData struct {
		ImageUrl        string
		Description     string
		Name            string
		MaxQuantity     int
		CurrentQuantity int
		Category        string
	}

	err = json.Unmarshal([]byte(itemDataJSON), &itemData)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	err = DeleteCld(c, itemData.ImageUrl)
	if err != nil {
		c.JSON(200, err)
		return
	}

	url, err := UploadToCld(c, file, fileHeader.Filename)
	if err != nil {
		c.JSON(200, err)
		return
	}
	itemData.ImageUrl = url

	if err := conect.DB.Model(&model.Item{}).Where("id=?", id).Updates(itemData).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to update item: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Item updated successfully"})
}

func PutItemByIDNoImage(c *gin.Context) {

	id := c.Param("id")
	itemDataJSON := c.PostForm("itemData")

	var itemData struct {
		Description     string
		Name            string
		MaxQuantity     int
		CurrentQuantity int
		Category        string
	}

	err := json.Unmarshal([]byte(itemDataJSON), &itemData)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		fmt.Println(err.Error())
		return
	}
	if err := conect.DB.Model(&model.Item{}).Where("id=?", id).Updates(itemData).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to update item: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Item updated successfully"})
}
