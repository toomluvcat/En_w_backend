package handler

import (
	"Render/app/conect"
	"Render/app/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func CreateUser(c *gin.Context) {
	var req model.User
	if err := c.BindJSON(&req); err != nil {
		c.Status(400)
		return
	}
	result := conect.DB.Select("Name","Email","StudentID","Major").Create(&req)

	if result.Error != nil {
		if pgErr, ok := result.Error.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(400, gin.H{"error": "Student ID OR Email have match another Student ID OR Email"})
			return
		}
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to insert User: %v", result.Error)})
		return
	}
	c.JSON(201, gin.H{"message": "Create user successfully", "req": req})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "Not found id"})
		return
	}

	var userRes model.User
	if err := conect.DB.Preload("Events.Loans.Item").First(&userRes, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to select value"})
		return
	}

	
	type LoanItem struct {
		ItemID   uint
		Category string
		Name string
		Quantity int
	}

	type EventResponse struct {
		EventID   uint
		CreatedAt time.Time
		ApprovedAt time.Time
		Status    string
		Loan      []LoanItem
	}

	type UserResponse struct {
		UserID uint
		Name   string
		Major  string
		StudentID string
		Email  string
		Event []EventResponse
	}

	
	var events []EventResponse

	for _,e := range userRes.Events{
		var loans []LoanItem
		
		for _,l := range e.Loans{
			loans = append(loans, LoanItem{
				ItemID: l.Item.ID,
				Name: l.Item.Name,
				Category: l.Item.Category,
				Quantity: l.Quantity,
			})
		}
		
		var ApprovedAt time.Time
		if e.Status =="Pending"{
			ApprovedAt =e.UpdatedAt
		}

		events = append(events, EventResponse{
			EventID: e.ID,
			CreatedAt: e.CreatedAt,
			ApprovedAt: ApprovedAt,
			Status: e.Status,
			Loan: loans,
		})
		
	}
	user:=UserResponse{
		UserID: userRes.ID,
		Name: userRes.Name,
		Email: userRes.Email,
		Major: userRes.Major,
		StudentID: userRes.StudentID,
		Event: events,
	}
	

	c.JSON(200, user)
}

func PutUserByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "Not found id"})
		return
	}

	type ReqPut struct {
		Name      string
		StudentID string
		Major     string
	}

	var req ReqPut
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	if err := conect.DB.Model(&model.User{}).Where("users.id=?", id).Updates(map[string]interface{}{
		"name":       req.Name,
		"student_id": req.StudentID,
		"major":      req.Major,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to update value :%v", err)})
		return
	}
	c.Status(201)
}
