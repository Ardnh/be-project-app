package dto

// ================ DTO ====================
type ProjectsDto struct {
	ID           string  `json:"id" validate:"required"`
	UserID       string  `json:"user_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
	CategoryName string  `json:"category_name" validate:"required"`
}

// ================ REQUEST DTO ====================
type CreateProjectsDto struct {
	UserID       string  `json:"user_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
	CategoryName string  `json:"category_name" validate:"required"`
}

type UpdateProjectsDto struct {
}

// ================ RESPONSE DTO ====================
