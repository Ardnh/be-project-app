package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CategoryID   uuid.UUID `gorm:"type:uuid"`
	Category     Category  `gorm:"foreignKey:CategoryID"` // Relasi ke Category
	UserID       uuid.UUID `gorm:"type:uuid"`
	User         User      `gorm:"foreignKey:UserID"` // Relasi ke User
	Name         string    `gorm:"type:varchar(100)"`
	Description  string    `gorm:"type:text"`
	Budget       int
	ProjectItems []ProjectItem `gorm:"foreignKey:ProjectID"` // Relasi ke ProjectItem
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type ProjectItem struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ProjectID  uuid.UUID `gorm:"type:uuid"`
	Project    Project   `gorm:"foreignKey:ProjectID"` // Relasi ke Project
	Name       string
	BudgetItem int
	Status     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
