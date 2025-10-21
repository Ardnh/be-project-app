package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type projectTodolistItemService struct {
	repo repositories.ProjectTodolistItemRepository
}

func NewProjectTodolistItemService(repo repositories.ProjectTodolistItemRepository) services.ProjectTodolistItemService {
	return &projectTodolistItemService{
		repo: repo,
	}
}

func (r *projectTodolistItemService) CreateProjectsTodolistItem(ctx context.Context, todoItem *dto.CreateProjectTodolistItemsDto) error
func (r *projectTodolistItemService) UpdateProjectsTodolistItem(ctx context.Context, id string, user *dto.UpdateProjectTodolistItemsDto) error
func (r *projectTodolistItemService) DeleteProjectsTodolistItem(ctx context.Context, id string) error
