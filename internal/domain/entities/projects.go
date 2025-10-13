package entities

import "time"

//	type Project struct {
//		ID           string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
//		UserID       string     `json:"user_id" gorm:"type:uuid;not null;index"`
//		Name         string     `json:"name" gorm:"type:varchar(255);not null"`
//		Budget       float64    `json:"budget" gorm:"type:decimal(15,2);not null"`
//		CategoryName string     `json:"category_name" gorm:"type:varchar(100);not null;index"`
//		CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
//		UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
//		DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index:idx_projects_deleted_at"`
//	}
type Projects struct {
	ID           string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       string  `gorm:"index" json:"user_id"`
	Name         string  `gorm:"size:200" json:"name"`
	Budget       float64 `json:"budget"`
	CategoryName string  `gorm:"size:100" json:"category_name"`

	// Relasi
	User            *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectExpenses []ProjectExpenses `gorm:"foreignKey:ProjectID" json:"project_expenses,omitempty"`
	// ProjectTodolists []ProjectTodolists `gorm:"foreignKey:ProjectID" json:"project_todolists,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

type GetProjectsParams struct {
	CategoryName string
	Search       string
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}

type ProjectWithTodolistAndExpenses struct {
	ID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       string     `gorm:"type:uuid;not null;index"`
	Name         string     `gorm:"type:varchar(255);not null"`
	Budget       float64    `gorm:"type:decimal(15,2);not null"`
	CategoryName string     `gorm:"type:varchar(100);not null;index"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `gorm:"index:idx_projects_deleted_at"`
	// ProjectExpensesWithItem []*ProjectExpensesWithItem `json:"project_expenses"`
}
