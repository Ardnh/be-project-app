package services

import (
	"context"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/application/mapper"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	repositories "github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
)

type projectsService struct {
	projectsRepo repositories.ProjectsRepository
}

func NewProjectsService(repo repositories.ProjectsRepository) services.ProjectsService {
	return &projectsService{
		projectsRepo: repo,
	}
}

func (r *projectsService) GetProjectsByUserID(ctx context.Context, id string) ([]*dto.ProjectsDto, error) {

	projects, err := r.projectsRepo.FindByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	projectsDto := mapper.ToProjectsDTO(projects)

	return projectsDto, nil
}

func (r *projectsService) GetProjectsByID(ctx context.Context, id string) (*dto.ProjectsDto, error) {

	return nil, nil
}
func (r *projectsService) GetAllUProjects(ctx context.Context) ([]*dto.ProjectsDto, error) {

	return nil, nil
}
func (r *projectsService) CreateProjects(ctx context.Context, user *dto.CreateProjectsDto) error {

	projectEntities := &entities.Project{
		UserID:       user.UserID,
		Name:         user.Name,
		Budget:       user.Budget,
		CategoryName: user.CategoryName,
	}

	err := r.projectsRepo.Create(ctx, projectEntities)
	if err != nil {
		return err
	}

	return nil
}
