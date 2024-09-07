package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex"`
	Password  string    `gorm:"type:varchar(100)"`
	Projects  []Project `gorm:"foreignKey:UserID"` // Relasi ke Project
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	User      User      `gorm:"foreignKey:UserID"`
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	Linkedin  string
	Twitter   string
}
