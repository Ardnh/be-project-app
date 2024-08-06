package model

import "gorm.io/gorm"

type ProjectItem struct {
	*gorm.Model
	ProjectID  uint
	Project    Project `gorm:"foreignKey:ProjectID"`
	Name       string
	BudgetItem int
	Status     bool
}
