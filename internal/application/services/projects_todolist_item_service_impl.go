package services

import (
	"context"
	"time"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"github.com/google/uuid"
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

	// Parse string UUID menjadi uuid.UUID
	projectTodolistID, err := uuid.Parse(todoItem.ProjectTodolistID)
	if err != nil {
		// Handle error jika string bukan format UUID yang valid
		return err // atau handle sesuai kebutuhan
	}

	req := &entities.ProjectTodolistItems{
		ProjectTodolistID: projectTodolistID,
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
	// Parse string UUID menjadi uuid.UUID
	projectTodolistID, err := uuid.Parse(todoItem.ProjectTodolistID)
	todoItemID, errTodoItemID := uuid.Parse(todoItem.ID)
	if err != nil {
		// Handle error jika string bukan format UUID yang valid
		return err // atau handle sesuai kebutuhan
	}

	if errTodoItemID != nil {
		return errTodoItemID
	}

	req := &entities.ProjectTodolistItems{
		ID:                todoItemID,
		ProjectTodolistID: projectTodolistID,
		Name:              todoItem.Name,
		CategoryName:      todoItem.CategoryName,
		IsCompleted:       todoItem.IsCompleted,
		UpdatedAt:         time.Now(),
	}

	errCreate := r.repo.CreateProjectTodolistItem(ctx, req)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (r *projectTodolistItemService) DeleteProjectsTodolistItem(ctx context.Context, id string) error {

	return r.DeleteProjectsTodolistItem(ctx, id)
}
