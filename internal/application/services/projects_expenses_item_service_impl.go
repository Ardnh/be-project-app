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

type projectExpensesItemService struct {
	repo repositories.ProjectExpensesItemRepository
}

func NewProjectExpensesItemService(repo repositories.ProjectExpensesItemRepository) services.ProjectsExpensesItemService {
	return &projectExpensesItemService{
		repo: repo,
	}
}

func (r *projectExpensesItemService) CreateProjectsExpensesItem(ctx context.Context, expensesItem *dto.CreateProjectExpensesItemDto) (*dto.ProjectExpensesItemDto, error) {

	req := entities.ProjectExpenseItems{
		ProjectExpenseID: expensesItem.ProjectExpensesId,
		Name:             expensesItem.Name,
		Amount:           expensesItem.Amount,
		CategoryName:     expensesItem.CategoryName,
	}

	createdItem, err := r.repo.Create(ctx, &req)
	if err != nil {
		return nil, err
	}

	expenseItemDto := mapper.ToProjectExpensesItemDTO(*createdItem)
	if expenseItemDto == nil {
		return nil, errors.New("Failed to convert expense item to DTO")
	}

	return expenseItemDto, nil
}

func (r *projectExpensesItemService) UpdateProjectsExpensesItem(ctx context.Context, id string, expensesItem *dto.UpdateProjectExpensesItemDto) (*dto.ProjectExpensesItemDto, error) {

	req := entities.ProjectExpenseItems{
		ID:               id,
		ProjectExpenseID: expensesItem.ProjectExpensesId,
		Name:             expensesItem.Name,
		Amount:           expensesItem.Amount,
		CategoryName:     expensesItem.CategoryName,
		UpdatedAt:        time.Now(),
	}

	updatedItem, err := r.repo.Update(ctx, &req)
	if err != nil {
		return nil, err
	}

	expenseItemDto := mapper.ToProjectExpensesItemDTO(*updatedItem)
	if expenseItemDto == nil {
		return nil, errors.New("Failed to convert expense item to DTO")
	}

	return expenseItemDto, nil
}

func (r *projectExpensesItemService) DeleteProjectsExpensesItem(ctx context.Context, id string) error {

	err := r.repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
