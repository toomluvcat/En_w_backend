package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	StudentID string
	Events    []Event `gorm:"foreignKey:UserID;references:ID"`
	Major     string
	Email     string `gorm:"unique"`
}

type Event struct {
	gorm.Model
	Status string
	UserID uint
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Loans  []Loan `gorm:"foreignKey:EventID;references:ID"`
}

type Item struct {
	gorm.Model
	Name            string `gorm:"unique"`
	MaxQuantity     int    `gorm:"check:max_quantity>=0"`
	CurrentQuantity int    `gorm:"check:current_quantity>=0"`
}

type Loan struct {
	gorm.Model
	Quantity int
	ItemID   uint
	EventID  uint
	Item     Item  `gorm:"foreignKey:ItemID;references:ID"`
	Event    Event `gorm:"constraint:OnDelete:CASCADE;"`
}
