package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type projectTodolistService struct {
	repo repositories.ProjectTodolistRepository
}

func NewProjectTodolistService(repo repositories.ProjectTodolistRepository) services.ProjectTodolistService {
	return &projectTodolistService{
		repo: repo,
	}
}

func (r *projectTodolistService) CreateProjectsTodolist(ctx context.Context, todo *dto.CreateProjectTodolistsDto) error
func (r *projectTodolistService) UpdateProjectsTodolist(ctx context.Context, id string, user *dto.UpdateProjectTodolistsDto) error
func (r *projectTodolistService) DeleteProjectsTodolist(ctx context.Context, id string) error
