package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type ProjectTodolistRepository interface {
	CreateProjectTodolist(ctx context.Context, todo *entities.ProjectTodolists) (*entities.ProjectTodolists, error)
	UpdateProjectTodolist(ctx context.Context, todo *entities.ProjectTodolists) (*entities.ProjectTodolists, error)
	DeleteProjectTodolist(ctx context.Context, id string) error
}
