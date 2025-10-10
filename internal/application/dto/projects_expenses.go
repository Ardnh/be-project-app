package dto

// ================ DTO ====================
type ProjectsExpensesDto struct {
	ID        string `json:"id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

// ================ REQUEST DTO ====================
type CreateProjectsExpensesDto struct {
	ID        string `json:"id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}
