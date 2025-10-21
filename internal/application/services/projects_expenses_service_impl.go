package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
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

func (s *projectExpensesService) CreateProjectsExpenses(ctx context.Context, expenses *dto.CreateProjectsExpensesDto) error {

	req := &entities.ProjectExpenses{
		ProjectID: expenses.ProjectID,
		Name:      expenses.Name,
	}

	err := s.projectsExpensesRepo.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *projectExpensesService) UpdateProjectsExpenses(ctx context.Context, id string, expenses *dto.UpdateProjectsExpensesDto) error {

	req := &entities.ProjectExpenses{
		ID:        id,
		ProjectID: expenses.ProjectID,
		Name:      expenses.Name,
	}

	err := s.projectsExpensesRepo.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *projectExpensesService) DeleteProjectsExpenses(ctx context.Context, id string) error {

	err := s.projectsExpensesRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
