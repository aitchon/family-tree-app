package models

import "gorm.io/gorm"

type UserEdit struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	TreeID   uint   `gorm:"not null"`
	EditData string `gorm:"not null"`          // JSON or other format for proposed changes
	Status   string `gorm:"default:'pending'"` // e.g., "pending", "approved", "rejected"
}

type ModerationQueue struct {
	gorm.Model
	EditID  uint   `gorm:"not null"`
	AdminID uint   `gorm:"not null"` // Admin who reviewed the edit
	Status  string `gorm:"not null"` // e.g., "approved", "rejected"
}
