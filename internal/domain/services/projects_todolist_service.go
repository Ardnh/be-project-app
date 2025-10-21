package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectTodolistService interface {
	CreateProjectsTodolist(ctx context.Context, todo *dto.CreateProjectTodolistsDto) error
	UpdateProjectsTodolist(ctx context.Context, id string, user *dto.UpdateProjectTodolistsDto) error
	DeleteProjectsTodolist(ctx context.Context, id string) error
}
