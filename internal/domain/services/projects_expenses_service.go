package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectsExpensesService interface {
	CreateProjectsExpenses(ctx context.Context, expenses *dto.CreateProjectsExpensesDto) (*dto.ProjectsExpensesDto, error)
	UpdateProjectsExpenses(ctx context.Context, id string, user *dto.UpdateProjectsExpensesDto) (*dto.ProjectsExpensesDto, error)
	DeleteProjectsExpenses(ctx context.Context, id string) error
}
