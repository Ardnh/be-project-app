package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type ProjectTodolistItemRepository interface {
	CreateProjectTodolistItem(ctx context.Context, todo *entities.ProjectTodolistItems) (*entities.ProjectTodolistItems, error)
	UpdateProjectTodolistItem(ctx context.Context, todo *entities.ProjectTodolistItems) (*entities.ProjectTodolistItems, error)
	DeleteProjectTodolistItem(ctx context.Context, id string) error
}
