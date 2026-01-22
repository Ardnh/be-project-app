package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectTodolistService interface {
	CreateProjectsTodolist(ctx context.Context, todo *dto.CreateProjectTodolistsDto) (*dto.ProjectTodolistsDto, error)
	UpdateProjectsTodolist(ctx context.Context, id string, user *dto.UpdateProjectTodolistsDto) (*dto.ProjectTodolistsDto, error)
	DeleteProjectsTodolist(ctx context.Context, id string) error
}
