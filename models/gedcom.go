package models

import "gorm.io/gorm"

type GEDCOMFile struct {
	gorm.Model
	UserID   uint   `gorm:"not null"` // ID of the user who uploaded the file
	Filename string `gorm:"not null"` // Name of the uploaded file
	Content  []byte `gorm:"not null"` // Raw content of the .ged file
}
