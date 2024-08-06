package schema

import "gorm.io/gorm"

type Users struct {
	*gorm.Model
	Username string
	Email    string
	Password string
}

type UserProfile struct {
	*gorm.Model
	UserID    int
	Bio       string
	Role      string
	Facebook  string
	Instagram string
	Linkedin  string
	Twitter   string
}

type FollowUsers struct {
	*gorm.Model
	UserID          int
	FollowingUserID int
}
