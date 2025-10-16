package services

import (
	"context"
	"time"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type projectExpensesItemService struct {
	repo repositories.ProjectExpensesItemRepository
}

func NewProjectExpensesItemService(repo repositories.ProjectExpensesItemRepository) services.ProjectsExpensesItemService {
	return &projectExpensesItemService{
		repo: repo,
	}
}

func (r *projectExpensesItemService) CreateProjectsExpensesItem(ctx context.Context, expensesItem *dto.CreateProjectExpensesItemDto) error {

	req := entities.ProjectExpenseItems{
		ProjectExpenseID: expensesItem.ProjectExpensesId,
		Name:             expensesItem.Name,
		Amount:           expensesItem.Amount,
		CategoryName:     expensesItem.CategoryName,
	}

	err := r.repo.Create(ctx, &req)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectExpensesItemService) UpdateProjectsExpensesItem(ctx context.Context, id string, expensesItem *dto.UpdateProjectExpensesItemDto) error {

	req := entities.ProjectExpenseItems{
		ID:               id,
		ProjectExpenseID: expensesItem.ProjectExpensesId,
		Name:             expensesItem.Name,
		Amount:           expensesItem.Amount,
		CategoryName:     expensesItem.CategoryName,
		UpdatedAt:        time.Now(),
	}

	err := r.repo.Update(ctx, &req)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectExpensesItemService) DeleteProjectsExpensesItem(ctx context.Context, id string) error {

	err := r.repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
