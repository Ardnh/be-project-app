package dto

// ================ DTO ====================
type ProjectTodolistItemsDto struct {
	ID                string  `json:"id"`
	ProjectTodolistID string  `json:"project_todolist_id"`
	Name              string  `json:"name"`
	CategoryName      string  `json:"category_name"`
	IsCompleted       bool    `json:"is_completed"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	DeletedAt         *string `json:"deleted_at"`
}

// ================ REQUEST DTO ====================
type CreateProjectTodolistItemsDto struct {
	ProjectTodolistID string `json:"project_todolist_id" validate:"required"`
	Name              string `json:"name" validate:"required"`
	CategoryName      string `json:"category_name" validate:"required"`
	IsCompleted       *bool  `json:"is_completed" validate:"required"`
}
type UpdateProjectTodolistItemsDto struct {
	ID           string `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
	IsCompleted  *bool  `json:"is_completed" validate:"required"`
}
