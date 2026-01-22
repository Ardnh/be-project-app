package services

import (
	"context"
	"errors"
	"time"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/application/mapper"
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

func (r *projectTodolistItemService) CreateProjectsTodolistItem(ctx context.Context, todoItem *dto.CreateProjectTodolistItemsDto) (*dto.ProjectTodolistItemsDto, error) {

	isCompletedStatus := false
	if todoItem.IsCompleted != nil {
		isCompletedStatus = *todoItem.IsCompleted
	}

	req := &entities.ProjectTodolistItems{
		ProjectTodolistID: todoItem.ProjectTodolistID,
		Name:              todoItem.Name,
		CategoryName:      todoItem.CategoryName,
		IsCompleted:       isCompletedStatus,
	}

	newTodo, errCreate := r.repo.CreateProjectTodolistItem(ctx, req)
	if errCreate != nil {
		return nil, errCreate
	}

	todoItemDto := mapper.ToProjectTodolistItemDTO(*newTodo)
	if todoItemDto == nil {
		return nil, errors.New("Failed to convert todo item to DTO")
	}

	return todoItemDto, nil
}

func (r *projectTodolistItemService) UpdateProjectsTodolistItem(ctx context.Context, id string, todoItem *dto.UpdateProjectTodolistItemsDto) (*dto.ProjectTodolistItemsDto, error) {

	isCompletedStatus := false
	if todoItem.IsCompleted != nil {
		isCompletedStatus = *todoItem.IsCompleted
	}

	req := &entities.ProjectTodolistItems{
		ID:           todoItem.ID,
		Name:         todoItem.Name,
		CategoryName: todoItem.CategoryName,
		IsCompleted:  isCompletedStatus,
		UpdatedAt:    time.Now(),
	}

	updatedTodoItem, errUpdate := r.repo.UpdateProjectTodolistItem(ctx, req)
	if errUpdate != nil {
		return nil, errUpdate
	}

	todoItemDto := mapper.ToProjectTodolistItemDTO(*updatedTodoItem)
	if todoItemDto == nil {
		return nil, errors.New("Failed to convert todo item to DTO")
	}

	return todoItemDto, nil
}

func (r *projectTodolistItemService) DeleteProjectsTodolistItem(ctx context.Context, id string) error {

	return r.repo.DeleteProjectTodolistItem(ctx, id)
}
