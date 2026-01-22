package services

import (
	"context"
	"errors"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/application/mapper"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"github.com/google/uuid"
)

type projectTodolistService struct {
	repo repositories.ProjectTodolistRepository
}

func NewProjectTodolistService(repo repositories.ProjectTodolistRepository) services.ProjectTodolistService {
	return &projectTodolistService{
		repo: repo,
	}
}

func (r *projectTodolistService) CreateProjectsTodolist(ctx context.Context, todo *dto.CreateProjectTodolistsDto) (*dto.ProjectTodolistsDto, error) {

	req := entities.ProjectTodolists{
		ProjectID: todo.ProjectID,
		Name:      todo.Name,
	}

	createdTodo, errCreate := r.repo.CreateProjectTodolist(ctx, &req)
	if errCreate != nil {
		return nil, errCreate
	}

	createdTodoDto := mapper.ToProjectTodolistDTO(createdTodo)
	if createdTodoDto == nil {
		return nil, errors.New("Failed to convert todo item to DTO")
	}

	return createdTodoDto, nil
}

func (r *projectTodolistService) UpdateProjectsTodolist(ctx context.Context, id string, todo *dto.UpdateProjectTodolistsDto) (*dto.ProjectTodolistsDto, error) {

	req := entities.ProjectTodolists{
		ID:        todo.ID,
		ProjectID: todo.ProjectID,
		Name:      todo.Name,
	}

	updatedTodo, errUpdate := r.repo.UpdateProjectTodolist(ctx, &req)
	if errUpdate != nil {
		return nil, errUpdate
	}

	updatedTodoDto := mapper.ToProjectTodolistDTO(updatedTodo)
	if updatedTodoDto == nil {
		return nil, errors.New("Failed to convert todo item to DTO")
	}

	return updatedTodoDto, nil
}

func (r *projectTodolistService) DeleteProjectsTodolist(ctx context.Context, id string) error {

	parseId, errParseId := uuid.Parse(id)
	if errParseId != nil {
		return errParseId
	}

	err := r.repo.DeleteProjectTodolist(ctx, parseId.String())
	if err != nil {
		return err
	}

	return nil
}
