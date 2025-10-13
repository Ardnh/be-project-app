package dto

// ================ DTO ====================
type ProjectsExpensesDto struct {
	ID        string `json:"id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}
type ProjectsExpensesWithItemDto struct {
	ID        string `json:"id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	// ProjectExpensesItem []*ProjectExpensesItemDto `json:"project_expenses_item" validate:"required"`
}

// ================ REQUEST DTO ====================
type CreateProjectsExpensesDto struct {
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

type UpdateProjectsExpensesDto struct {
	ID        string `json:"id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}
