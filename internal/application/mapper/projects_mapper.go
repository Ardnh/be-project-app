package mapper

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

func ToProjectsDTO(projects []*entities.Project) []*dto.ProjectsDto {
	projectsDto := make([]*dto.ProjectsDto, 0, len(projects))

	for _, project := range projects {
		projectsDto = append(projectsDto, ToProjectDTO(project))
	}

	return projectsDto
}

func ToProjectDTO(project *entities.Project) *dto.ProjectsDto {
	if project == nil {
		return nil
	}

	return &dto.ProjectsDto{
		ID:           project.ID,
		UserID:       project.UserID,
		Name:         project.Name,
		Budget:       project.Budget,
		CategoryName: project.CategoryName,
	}
}
