package repositories

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

type ProjectsRepository interface {
	FindByUserID(ctx context.Context, userId string, params entities.GetProjectsParams) ([]*entities.Projects, error)
	FindProjectCategoryByUserID(ctx context.Context, userId string) ([]*entities.ProjectCategorySummary, error)
	FindProjectSummaryByUserID(ctx context.Context, userId string) (*entities.ProjectUserSummary, error)
	FindAll(ctx context.Context, params entities.GetProjectsParams) ([]*entities.Projects, error)
	FindByProjectId(ctx context.Context, id string) (*entities.Projects, error)
	Create(ctx context.Context, user *entities.Projects) error
	Update(ctx context.Context, user *entities.Projects) error
	Delete(ctx context.Context, id string) error
}
