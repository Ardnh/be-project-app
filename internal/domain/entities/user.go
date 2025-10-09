package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string     `gorm:"type:varchar(100);uniqueIndex:idx_users_username;not null"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex:idx_users_email;not null"`
	Password  string     `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `gorm:"index:idx_users_deleted_at"` // Pointer karena nullable
}
