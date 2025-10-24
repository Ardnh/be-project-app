package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectsService interface {
	GetProjectsByUserID(ctx context.Context, userId string, params dto.GetProjectsParams) ([]*dto.ProjectsDto, error)
	GetProjectsSummaryByUserID(ctx context.Context, userId string) (*dto.ProjectSummaryDto, error)
	GetProjectsByID(ctx context.Context, id string) (*dto.ProjectWithTodolistAndExpensesDto, error)
	GetProjectCategoryByUserId(ctx context.Context, userId string) ([]*dto.ProjectCategorySummaryDto, error)
	GetAllProjects(ctx context.Context, params dto.GetProjectsParams) ([]*dto.ProjectsDto, error)
	CreateProjects(ctx context.Context, user *dto.CreateProjectDto) error
	UpdateProjects(ctx context.Context, id string, user *dto.UpdateProjectDto) error
	DeleteProjects(ctx context.Context, id string) error
}
