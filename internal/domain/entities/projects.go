package entities

import "time"

type Project struct {
	ID           string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       string     `json:"user_id" gorm:"type:uuid;not null;index"`
	Name         string     `json:"name" gorm:"type:varchar(255);not null"`
	Budget       float64    `json:"budget" gorm:"type:decimal(15,2);not null"`
	CategoryName string     `json:"category_name" gorm:"type:varchar(100);not null;index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index:idx_projects_deleted_at"`
}

type GetProjectsParams struct {
	CategoryName string
	Search       string
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}
