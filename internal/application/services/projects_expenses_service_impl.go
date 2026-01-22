package services

import (
	"context"
	"errors"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/application/mapper"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type projectExpensesService struct {
	projectsExpensesRepo repositories.ProjectsExpensesRepository
}

func NewProjectExpensesService(repo repositories.ProjectsExpensesRepository) services.ProjectsExpensesService {
	return &projectExpensesService{
		projectsExpensesRepo: repo,
	}
}

func (s *projectExpensesService) CreateProjectsExpenses(ctx context.Context, expenses *dto.CreateProjectsExpensesDto) (*dto.ProjectsExpensesDto, error) {

	req := &entities.ProjectExpenses{
		ProjectID: expenses.ProjectID,
		Name:      expenses.Name,
	}

	createdExpenses, err := s.projectsExpensesRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	createdExpensesDto := mapper.ToProjectExpensesDTO(createdExpenses)
	if createdExpenses == nil {
		return nil, errors.New("Failed to convert expenses to DTO")
	}

	return createdExpensesDto, nil
}

func (s *projectExpensesService) UpdateProjectsExpenses(ctx context.Context, id string, expenses *dto.UpdateProjectsExpensesDto) (*dto.ProjectsExpensesDto, error) {

	req := &entities.ProjectExpenses{
		ID:   id,
		Name: expenses.Name,
	}

	updatedExpense, err := s.projectsExpensesRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	updatedExpenseDto := mapper.ToProjectExpensesDTO(updatedExpense)
	if updatedExpense == nil {
		return nil, errors.New("Failed to convert expenses to DTO")
	}

	return updatedExpenseDto, nil
}

func (s *projectExpensesService) DeleteProjectsExpenses(ctx context.Context, id string) error {

	err := s.projectsExpensesRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
