package schema

import (
	"github.com/google/uuid"
)

type Projects struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CategoryID   uuid.UUID
	Category     Categories `gorm:"foreignKey:CategoryID"`
	UserID       uuid.UUID
	Users        Users  `gorm:"foreignKey:UserID"`
	Name         string `gorm:"type:varchar(100)"`
	Description  string `gorm:"type:text"`
	Budget       int
	ProjectItems []ProjectItem
}
