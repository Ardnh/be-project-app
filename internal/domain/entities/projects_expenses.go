package entities

import "time"

type ProjectExpenses struct {
	ID        string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectID string `gorm:"type:uuid;index;not null" json:"project_id"`
	Name      string `json:"name" gorm:"type:varchar(255);not null"`

	Project            *Projects             `gorm:"-" json:"-"`
	ProjectExpenseItem []ProjectExpenseItems `gorm:"foreignKey:ProjectExpenseID" json:"project_expense_items,omitempty"`

	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index:idx_projects_deleted_at"`
}
