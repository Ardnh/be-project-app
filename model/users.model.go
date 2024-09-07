package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string     `json:"id"` // Ubah dari uint ke string jika menggunakan UUID
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type IsEmail struct {
	Email string `validate:"required,email"`
}

type Profile struct {
	*gorm.Model
	UserId    int
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	LinkedIn  string
	Twitter   string
}

type UserWithProfile struct {
	UserId         int
	FollowedUserId int
	Role           string
	Username       string
}

type ProfileUpdateRequestBody struct {
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	LinkedIn  string
	Twitter   string
}

type ProfileUpdateRequest struct {
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	Linkedin  string
	Twitter   string
}

type ProfileCreateRequest struct {
	UserId    string
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	Linkedin  string
	Twitter   string
}
