package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectTodolistItemService interface {
	CreateProjectsTodolistItem(ctx context.Context, todoItem *dto.CreateProjectTodolistItemsDto) error
	UpdateProjectsTodolistItem(ctx context.Context, id string, user *dto.UpdateProjectTodolistItemsDto) error
	DeleteProjectsTodolistItem(ctx context.Context, id string) error
}
