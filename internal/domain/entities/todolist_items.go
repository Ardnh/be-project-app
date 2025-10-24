package entities

import (
	"time"

	"gorm.io/gorm"
)

type ProjectTodolistItems struct {
	ID                string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectTodolistID string         `gorm:"type:uuid;not null;index"`
	Name              string         `gorm:"type:varchar(200);not null"`
	CategoryName      string         `gorm:"type:varchar(200)"`
	IsCompleted       bool           `gorm:"default:false"` // Lebih jelas dari "Status"
	CreatedAt         time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
