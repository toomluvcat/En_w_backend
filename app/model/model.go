package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string
	StudentID      string
	Major          string
	Email          string  `gorm:"unique"`
	BookmarksItems []Item  `gorm:"many2many:bookmarks"`
	Events         []Event `gorm:"foreignKey:UserID;references:ID"`
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
	Description     string
	Category        string
	MaxQuantity     int    `gorm:"check:max_quantity>=0"`
	CurrentQuantity int    `gorm:"check:current_quantity>=0"`
	BookmarkedBy    []User `gorm:"many2many:bookmarks;"`
	ImageUrl        string
}

type Loan struct {
	gorm.Model
	Quantity int
	ItemID   uint
	EventID  uint
	Item     Item  `gorm:"foreignKey:ItemID;references:ID"`
	Event    Event `gorm:"constraint:OnDelete:CASCADE;"`
}
