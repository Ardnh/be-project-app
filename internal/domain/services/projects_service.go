package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
)

type ProjectsService interface {
	GetProjectsByUserID(ctx context.Context, userId string) ([]*dto.ProjectsDto, error)
	GetProjectsByID(ctx context.Context, id string) (*dto.ProjectsDto, error)
	GetAllUProjects(ctx context.Context) ([]*dto.ProjectsDto, error)
	CreateProjects(ctx context.Context, user *dto.CreateProjectsDto) error
}
