package model

import (
	"gorm.io/gorm"
)

type Project struct {
	*gorm.Model
	CategoryID   uint
	Category     Category `gorm:"foreignKey:CategoryID"`
	Name         string   `gorm:"type:varchar(100)"`
	Description  string   `gorm:"type:text"`
	Budget       int
	ProjectItems []ProjectItem
}
