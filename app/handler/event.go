package handler

import (
	"Render/app/conect"
	"Render/app/model"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllEvent(c *gin.Context) {

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
	if err := conect.DB.Preload("User").Preload("Loans.Item").Find(&event).Error; err != nil {
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

	c.JSON(200, res)
}

func GetEventByUserID(c *gin.Context) {
	ID := c.Param("id")

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
	if err := conect.DB.Preload("User").Preload("Loans.Item").Where("user_id = ?", ID).Find(&event).Error; err != nil {
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

	c.JSON(200, res)
}

func CreateEvent(c *gin.Context) {
	type Req struct {
		UserID uint
		Items  []struct {
			ItemID   uint
			Quantity int
		}
	}

	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	var user model.User
	if err := conect.DB.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "Student not found"})
		return
	}

	tx := conect.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	event := model.Event{
		Status: "Pending",
		UserID: req.UserID,
	}

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Fail to create event" + err.Error()})
		return
	}

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			tx.Rollback()
			c.JSON(400, gin.H{"error": fmt.Sprintf("Quantity must be greater than 0 for item ID %d", itemReq.ItemID)})
			return
		}

		var item model.Item
		if err := tx.Where("id = ?", itemReq.ItemID).First(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(404, gin.H{"error": "item id not match with any item id"})
			return
		}

		loan := model.Loan{

			Quantity: itemReq.Quantity,
			EventID:  event.ID,
			ItemID:   itemReq.ItemID,
		}

		if err := tx.Create(&loan).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Fail to create loan: " + err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to commit transaction"})
		return
	}

	c.JSON(200, gin.H{"message": "Loans create successfully"})

}

func PutEventdByID(c *gin.Context) {
	id := c.Param("id")

	type Req struct {
		Status string
	}

	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	var event model.Event
	err := conect.DB.Preload("Loans").First(&event, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch event: " + err.Error(),
		})
		return
	}

	tx := conect.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Status == "complete" {

		result := tx.Model(&model.Event{}).Where("id = ? AND status = ? ", id, "approved").Update("status", req.Status)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Fail to update Event Status"})
			return
		}

		caseStmt := "UPDATE items Set current_quantity = current_quantity + CASE id"
		ids := []string{}
		for _, loan := range event.Loans {
			caseStmt += fmt.Sprintf(" WHEN %d THEN %d", loan.ItemID, loan.Quantity)
			ids = append(ids, fmt.Sprintf("%d", loan.ItemID))
		}
		caseStmt += fmt.Sprintf(" END WHERE id IN (%s)", strings.Join(ids, ","))
		fmt.Println(caseStmt)
		if err := tx.Exec(caseStmt).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Fail to update items" + err.Error()})
			return
		}

	} else if req.Status == "reject" {
		result := tx.Model(&model.Event{}).Where("id = ? AND status = ? ", id, "pending").Update("status", req.Status)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Fail to update Event Status"})
			return
		}
		c.JSON(200, gin.H{"error": fmt.Sprintf("reject successfully event id: %d", id)})
		return

	} else {
		result := tx.Model(&model.Event{}).Where("id = ? AND status = ? ", id, "pending").Update("status", req.Status)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Fail to update Event Status"})
			return
		}
		for _, loan := range event.Loans {
			var item model.Item

			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&item, loan.ItemID).Error; err != nil {
				tx.Rollback()
				c.JSON(400, gin.H{"error": fmt.Sprintf("Item not found at id: %d", loan.ItemID)})
				return
			}

			if item.CurrentQuantity < loan.Quantity {
				tx.Rollback()
				c.JSON(400, gin.H{"error": fmt.Sprintf("Item not enough quantity for item: %d", loan.ItemID)})
				return
			}

			item.CurrentQuantity -= loan.Quantity
			if err := tx.Save(item).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Fail to update item quantity:" + err.Error()})
				return
			}
		}

	}
	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": "Fail to commit transaction: " + err.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "status change successfully as id: " + id})
}

func DeleteItemByID(c *gin.Context) {
	id := c.Param("id")

	if result := conect.DB.Delete(&model.Event{}, id).Error; result != nil {
		c.Status(500)
		return
	}

	c.Status(200)
}
