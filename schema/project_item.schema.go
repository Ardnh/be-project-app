package schema

import "github.com/google/uuid"

type ProjectItem struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ProjectID  uint
	Project    Projects `gorm:"foreignKey:ProjectID"`
	Name       string
	BudgetItem int
	Status     bool
}
