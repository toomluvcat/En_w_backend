package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	StudentID string  `gorm:"unique"`
	Events    []Event `gorm:"foreignKey:StudentID;references:StudentID"`
	Major     string
	Email     string `gorm:"unique"`
}

type Event struct {
	gorm.Model
	StudentID string
	Time      time.Time
	Loans     []Loan `gorm:"foreignKey:ItemID;references:ID"`
}

type Item struct {
	gorm.Model
	Name            string `gorm:"unique"`
	MaxQuantity     int
	CurrentQuantity int
}

type Loan struct {
	gorm.Model
	Quantity int
	ItemID   uint
	EventID uint
	Item     Item `gorm:"foreignKey:ItemID;references:ID"`
}
