package dto

// ================ DTO ====================
type ProjectsExpensesDto struct {
	ID           string                    `json:"id" validate:"required"`
	ProjectID    string                    `json:"project_id" validate:"required"`
	Name         string                    `json:"name" validate:"required"`
	ExpensesUsed float64                   `json:"expenses_used" validate:"required"`
	ExpensesItem []*ProjectExpensesItemDto `json:"expenses_item" validate:"required"`
}

// ================ REQUEST DTO ====================
type CreateProjectsExpensesDto struct {
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

type UpdateProjectsExpensesDto struct {
	Name string `json:"name" validate:"required"`
}
