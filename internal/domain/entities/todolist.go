package entities

import (
	"time"

	"gorm.io/gorm"
)

type ProjectTodolists struct {
	ID        string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectID string `gorm:"type:uuid;not null;index"`
	Name      string `gorm:"type:varchar(200);not null"`

	ProjectTodolistItems []ProjectTodolistItems `gorm:"foreignKey:ProjectTodolistID"`

	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
