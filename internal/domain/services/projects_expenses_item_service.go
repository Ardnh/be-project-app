package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectsExpensesItemService interface {
	CreateProjectsExpensesItem(ctx context.Context, expensesItem *dto.CreateProjectExpensesItemDto) error
	UpdateProjectsExpensesItem(ctx context.Context, id string, user *dto.UpdateProjectExpensesItemDto) error
	DeleteProjectsExpensesItem(ctx context.Context, id string) error
}
