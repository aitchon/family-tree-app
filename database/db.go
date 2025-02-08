package database

import (
	"family-tree-app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("family-tree.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto-migrate all models
	DB.AutoMigrate(&models.User{}, &models.GEDCOMFile{}, &models.FamilyTree{}, &models.UserEdit{}, &models.ModerationQueue{})
}
