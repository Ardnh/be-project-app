package model

import (
	"time"

	"github.com/google/uuid"
)

type Projects struct {
	ID          uuid.UUID  `json:"id"`
	CategoryID  uuid.UUID  `json:"category_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Budget      int        `json:"budget"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Category    Category   `json:"category" gorm:"foreignKey:CategoryID"`
}

type ProjectItem struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Name       string    `json:"name"`
	BudgetItem int       `json:"budget_item"`
	Status     bool      `json:"status"`
}

type CreateProjectRequest struct {
	CategoryID  uuid.UUID `json:"category_id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Budget      int       `json:"budget"`
}

type UpdateProjectRequest struct {
	CategoryID  uuid.UUID `json:"category_id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Budget      int       `json:"budget"`
}

type CreateProjectItem struct {
	ProjectID  uuid.UUID `json:"project_id"`
	Name       string    `json:"name"`
	BudgetItem int       `json:"budget_item"`
	Status     bool      `json:"status"`
}

type UpdateProjectItem struct {
	ProjectID  uuid.UUID `json:"project_id"`
	Name       string    `json:"name"`
	BudgetItem int       `json:"budget_item"`
	Status     bool      `json:"status"`
}

type ProjectItemResponse struct {
	ProjectItems []ProjectItem `json:"project_items"`
	Project      Projects      `json:"projects"`
}
