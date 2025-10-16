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
	// Length 0, capacity 3
	expenses := make([]*dto.ProjectsExpensesDto, 0, len(project.ProjectExpenses))

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

	expensesItem := make([]*dto.ProjectExpensesItemDto, 0, len(expenses.ProjectExpenseItem))

	for _, item := range expenses.ProjectExpenseItem {
		expensesItem = append(expensesItem, ToProjectExpensesItemDTO(item))
	}

	return &dto.ProjectsExpensesDto{
		ID:           expenses.ID,
		ProjectID:    expenses.ProjectID,
		Name:         expenses.Name,
		ExpensesItem: expensesItem,
	}
}

func ToProjectExpensesItemDTO(expensesItem entities.ProjectExpenseItems) *dto.ProjectExpensesItemDto {
	return &dto.ProjectExpensesItemDto{
		ID:                expensesItem.ID,
		ProjectExpensesId: expensesItem.ProjectExpenseID,
		Name:              expensesItem.Name,
		Amount:            expensesItem.Amount,
		CategoryName:      expensesItem.CategoryName,
	}
}

func ToProjectCategorySummaryDTO(projectCategory []*entities.ProjectCategorySummary) []*dto.ProjectCategorySummaryDto {

	var result = make([]*dto.ProjectCategorySummaryDto, 0, len(projectCategory))
	for _, item := range projectCategory {
		result = append(result, &dto.ProjectCategorySummaryDto{
			CategoryName: item.CategoryName,
			Total:        item.Total,
		})
	}

	return result
}
