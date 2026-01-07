package services

import (
	"context"
	"time"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
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

func (r *projectTodolistItemService) CreateProjectsTodolistItem(ctx context.Context, todoItem *dto.CreateProjectTodolistItemsDto) error {

	req := &entities.ProjectTodolistItems{
		ProjectTodolistID: todoItem.ProjectTodolistID,
		Name:              todoItem.Name,
		CategoryName:      todoItem.CategoryName,
		IsCompleted:       todoItem.IsCompleted,
	}

	errCreate := r.repo.CreateProjectTodolistItem(ctx, req)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (r *projectTodolistItemService) UpdateProjectsTodolistItem(ctx context.Context, id string, todoItem *dto.UpdateProjectTodolistItemsDto) error {

	req := &entities.ProjectTodolistItems{
		ID:           todoItem.ID,
		Name:         todoItem.Name,
		CategoryName: todoItem.CategoryName,
		IsCompleted:  todoItem.IsCompleted,
		UpdatedAt:    time.Now(),
	}

	errCreate := r.repo.UpdateProjectTodolistItem(ctx, req)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (r *projectTodolistItemService) DeleteProjectsTodolistItem(ctx context.Context, id string) error {

	return r.repo.DeleteProjectTodolistItem(ctx, id)
}
