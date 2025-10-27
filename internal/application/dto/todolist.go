package dto

// ================ DTO ====================
type ProjectTodolistsDto struct {
	ID                   string                     `json:"id"`
	ProjectID            string                     `json:"project_id"`
	Name                 string                     `json:"name"`
	ProjectTodolistItems []*ProjectTodolistItemsDto `json:"todolist_items"`
	TotalTodo            int                        `json:"total_todo"`
	TotalCompletedTodo   int                        `json:"total_completed_todo"`
	IsTodolistCompleted  bool                       `json:"is_todolist_completed"`
	CreatedAt            string                     `json:"created_at"`
	UpdatedAt            string                     `json:"updated_at"`
	DeletedAt            *string                    `json:"deleted_at"`
}

// ================ REQUEST DTO ====================
type CreateProjectTodolistsDto struct {
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

type UpdateProjectTodolistsDto struct {
	ID        string `json:"id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	UpdatedAt string `json:"updated_at"`
}
