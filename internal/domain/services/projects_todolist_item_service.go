package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectTodolistItemService interface {
	CreateProjectsTodolistItem(ctx context.Context, todoItem *dto.CreateProjectTodolistItemsDto) (*dto.ProjectTodolistItemsDto, error)
	UpdateProjectsTodolistItem(ctx context.Context, id string, todoItem *dto.UpdateProjectTodolistItemsDto) (*dto.ProjectTodolistItemsDto, error)
	DeleteProjectsTodolistItem(ctx context.Context, id string) error
}
