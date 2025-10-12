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

func (r *projectsService) GetProjectsByUserID(ctx context.Context, userId string, params dto.GetProjectsParams) ([]*dto.ProjectsDto, error) {

	paramsEntities := entities.GetProjectsParams{
		CategoryName: params.CategoryName,
		Search:       params.Search,
		Limit:        params.Limit,
		Offset:       params.Offset,
		SortBy:       params.SortBy,
		SortOrder:    params.SortOrder,
	}

	projects, err := r.projectsRepo.FindByUserID(ctx, userId, paramsEntities)
	if err != nil {
		return nil, err
	}

	projectsDto := mapper.ToProjectsDTO(projects)

	return projectsDto, nil
}

func (r *projectsService) GetProjectsByID(ctx context.Context, id string) (*dto.ProjectsDto, error) {

	projectEntities, err := r.projectsRepo.FindByProjectId(ctx, id)
	if err != nil {
		return nil, err
	}

	projectDto := mapper.ToProjectDTO(projectEntities)

	return projectDto, nil
}

func (r *projectsService) GetAllProjects(ctx context.Context, params dto.GetProjectsParams) ([]*dto.ProjectsDto, error) {

	paramsEntities := entities.GetProjectsParams{
		CategoryName: params.CategoryName,
		Search:       params.Search,
		Limit:        params.Limit,
		Offset:       params.Offset,
		SortBy:       params.SortBy,
		SortOrder:    params.SortOrder,
	}

	projects, err := r.projectsRepo.FindAll(ctx, paramsEntities)
	if err != nil {
		return nil, err
	}

	projectsDto := mapper.ToProjectsDTO(projects)
	return projectsDto, nil
}

func (r *projectsService) CreateProjects(ctx context.Context, project *dto.CreateProjectsDto) error {

	projectEntities := &entities.Project{
		UserID:       project.UserID,
		Name:         project.Name,
		Budget:       project.Budget,
		CategoryName: project.CategoryName,
	}

	err := r.projectsRepo.Create(ctx, projectEntities)
	if err != nil {
		return err
	}

	return nil
}

func (r *projectsService) UpdateProjects(ctx context.Context, id string, project *dto.UpdateProjectsDto) error {

	projectEntities := &entities.Project{
		ID:           id,
		UserID:       project.UserID,
		Name:         project.Name,
		Budget:       project.Budget,
		CategoryName: project.CategoryName,
	}

	err := r.projectsRepo.Update(ctx, projectEntities)
	if err != nil {
		return err
	}

	return nil
}

func (r *projectsService) DeleteProjects(ctx context.Context, id string) error {

	err := r.projectsRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
