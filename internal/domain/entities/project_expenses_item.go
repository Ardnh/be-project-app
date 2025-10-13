package entities

import "time"

type ProjectExpenseItem struct {
	ID                string  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectExpensesID string  `json:"project_expenses_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string  `json:"name" gorm:"type:varchar(255);not null"`
	Amount            float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	CategoryName      string  `json:"category_name" gorm:"type:varchar(255);not null"`

	// Relasi
	ProjectExpense *ProjectExpenses `gorm:"foreignKey:ProjectExpensesID" json:"project_expense,omitempty"`

	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index:idx_projects_deleted_at"`
}
