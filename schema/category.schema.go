package schema

import "github.com/google/uuid"

type Categories struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name string
}
