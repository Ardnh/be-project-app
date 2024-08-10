package schema

import (
	"github.com/google/uuid"
)

type Users struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username string
	Email    string
	Password string
	Projects []Projects
}

type UserProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID
	User      Users `gorm:"foreignKey:UserID"`
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	Linkedin  string
	Twitter   string
}
