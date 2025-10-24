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
	todolists := make([]*dto.ProjectTodolistsDto, 0, len(project.ProjectTodolists))

	var usedBudget float64 = 0
	var totalTodolistItem int = 0
	var totalTodolistCompletedItem int = 0

	for _, expensesItem := range project.ProjectExpenses {
		item := ToProjectExpensesDTO(&expensesItem)
		expenses = append(expenses, item)

		usedBudget += item.ExpensesUsed
	}

	for _, todolist := range project.ProjectTodolists {
		todo := ToProjectTodolistDTO(&todolist)
		todolists = append(todolists, ToProjectTodolistDTO(&todolist))

		totalTodolistItem += todo.TotalTodo
		totalTodolistCompletedItem += todo.TotalCompletedTodo
	}

	return &dto.ProjectWithTodolistAndExpensesDto{
		ID:                         project.ID,
		UserID:                     project.UserID,
		Name:                       project.Name,
		Budget:                     project.Budget,
		StartDate:                  project.StartDate,
		EndDate:                    project.EndDate,
		BudgetUsed:                 usedBudget,
		CategoryName:               project.CategoryName,
		TotalTodolistItem:          totalTodolistItem,
		TotalTodolistCompletedItem: totalTodolistCompletedItem,
		ProjectExpenses:            expenses,
		ProjectTodolist:            todolists,
	}
}

func ToProjectExpensesDTO(expenses *entities.ProjectExpenses) *dto.ProjectsExpensesDto {

	expensesItem := make([]*dto.ProjectExpensesItemDto, 0, len(expenses.ProjectExpenseItem))
	var expensesUsed float64 = 0

	for _, item := range expenses.ProjectExpenseItem {
		itemDTO := ToProjectExpensesItemDTO(item)
		expensesItem = append(expensesItem, itemDTO)
		expensesUsed += itemDTO.Amount
	}

	return &dto.ProjectsExpensesDto{
		ID:           expenses.ID,
		ProjectID:    expenses.ProjectID,
		Name:         expenses.Name,
		ExpensesUsed: expensesUsed,
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

func ToProjectTodolistDTO(todolist *entities.ProjectTodolists) *dto.ProjectTodolistsDto {

	todolistItem := make([]*dto.ProjectTodolistItemsDto, 0, len(todolist.ProjectTodolistItems))

	for _, todoItem := range todolist.ProjectTodolistItems {
		var item = ToProjectTodolistItemDTO(todoItem)
		todolistItem = append(todolistItem, item)
	}

	totalTodo := len(todolist.ProjectTodolistItems)
	totalCompletedTodo := 0

	for _, todoItem := range todolist.ProjectTodolistItems {
		if todoItem.IsCompleted {
			totalCompletedTodo += 1
		}
	}

	return &dto.ProjectTodolistsDto{
		ID:                   todolist.ID,
		ProjectID:            todolist.ProjectID,
		Name:                 todolist.Name,
		TotalTodo:            totalTodo,
		TotalCompletedTodo:   totalCompletedTodo,
		ProjectTodolistItems: todolistItem,
	}
}

func ToProjectTodolistItemDTO(expensesItem entities.ProjectTodolistItems) *dto.ProjectTodolistItemsDto {
	return &dto.ProjectTodolistItemsDto{
		ID:                expensesItem.ID,
		ProjectTodolistID: expensesItem.ProjectTodolistID,
		Name:              expensesItem.Name,
		CategoryName:      expensesItem.CategoryName,
		IsCompleted:       expensesItem.IsCompleted,
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

func ToProjectSummaryDTO(summary *entities.ProjectUserSummary) *dto.ProjectSummaryDto {

	return &dto.ProjectSummaryDto{
		TotalProjects:          summary.TotalProject,
		TotalBudgetUsed:        summary.TotalBudget,
		TotalCompletedProjects: summary.TotalCompletedProject,
	}
}
