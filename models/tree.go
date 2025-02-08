package models

import "gorm.io/gorm"

type FamilyTree struct {
	gorm.Model
	GEDCOMFileID uint       `gorm:"not null"`
	GEDCOMFile   GEDCOMFile `gorm:"foreignKey:GEDCOMFileID"` // Belongs to a GEDCOM file
	Data         string     `gorm:"type:text;not null"`
	Edits        []UserEdit `gorm:"foreignKey:TreeID"` // A FamilyTree can have multiple UserEdits
}
