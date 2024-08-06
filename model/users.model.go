package model

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Username string
	Email    string
	Password string
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
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
	UserId    uint
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	Linkedin  string
	Twitter   string
}
