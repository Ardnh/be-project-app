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

	isCompleted := false
	if todoItem.IsCompleted != nil {
		isCompleted = *todoItem.IsCompleted
	}

	req := &entities.ProjectTodolistItems{
		ProjectTodolistID: todoItem.ProjectTodolistID,
		Name:              todoItem.Name,
		CategoryName:      todoItem.CategoryName,
		IsCompleted:       isCompleted,
	}

	errCreate := r.repo.CreateProjectTodolistItem(ctx, req)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (r *projectTodolistItemService) UpdateProjectsTodolistItem(ctx context.Context, id string, todoItem *dto.UpdateProjectTodolistItemsDto) error {

	isCompleted := false
	if todoItem.IsCompleted != nil {
		isCompleted = *todoItem.IsCompleted
	}

	req := &entities.ProjectTodolistItems{
		ID:                todoItem.ID,
		ProjectTodolistID: todoItem.ProjectTodolistID,
		Name:              todoItem.Name,
		CategoryName:      todoItem.CategoryName,
		IsCompleted:       isCompleted,
		UpdatedAt:         time.Now(),
	}

	errCreate := r.repo.CreateProjectTodolistItem(ctx, req)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (r *projectTodolistItemService) DeleteProjectsTodolistItem(ctx context.Context, id string) error {

	return r.repo.DeleteProjectTodolistItem(ctx, id)
}
