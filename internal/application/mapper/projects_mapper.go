package mapper

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

func ToProjectsDTO(projects []*entities.Projects) []*dto.ProjectsDto {
	projectsDto := make([]*dto.ProjectsDto, 0, len(projects))

	for _, project := range projects {
		projectsDto = append(projectsDto, ToProjectDTO(project))
	}

	return projectsDto
}

func ToProjectDTO(project *entities.Projects) *dto.ProjectsDto {
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

func ToProjectWithTodolistAndExpensesDTO(project *entities.Projects) *dto.ProjectWithTodolistAndExpensesDto {

	expenses := make([]*dto.ProjectsExpensesDto, len(project.ProjectExpenses))

	for _, expensesItem := range project.ProjectExpenses {
		expenses = append(expenses, ToProjectExpensesDTO(&expensesItem))
	}

	return &dto.ProjectWithTodolistAndExpensesDto{
		ID:              project.ID,
		UserID:          project.UserID,
		Name:            project.Name,
		Budget:          project.Budget,
		CategoryName:    project.CategoryName,
		ProjectExpenses: expenses,
	}
}

func ToProjectExpensesDTO(expenses *entities.ProjectExpenses) *dto.ProjectsExpensesDto {

	return &dto.ProjectsExpensesDto{
		ID:        expenses.ID,
		ProjectID: expenses.ProjectID,
		Name:      expenses.Name,
	}
}
