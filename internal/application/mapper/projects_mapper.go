package mapper

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	date_utils "github.com/Ardnh/be-project-app/internal/utils/date"
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
		ID:                project.ID,
		UserID:            project.UserID,
		Name:              project.Name,
		Budget:            project.Budget,
		CategoryName:      project.CategoryName,
		TotalTodolist:     0,
		TotalTodolistDone: 0,
	}
}

func ToProjectsByUserDTO(entities []*entities.ProjectsByUserId) []*dto.ProjectsByUserIdDto {
	if entities == nil {
		return []*dto.ProjectsByUserIdDto{}
	}

	dtos := make([]*dto.ProjectsByUserIdDto, len(entities))

	for i, entity := range entities {
		if entity == nil {
			continue
		}

		daysRemaining := date_utils.CalculateDaysRemaining(entity.EndDate)
		daysRemainingStatus := date_utils.GetDaysRemainingStatus(daysRemaining)
		dtos[i] = &dto.ProjectsByUserIdDto{
			ProjectID:             entity.ProjectID,
			UserID:                entity.UserID,
			Name:                  entity.Name,
			Budget:                entity.Budget,
			IsCompleted:           entity.IsCompleted,
			CategoryName:          entity.CategoryName,
			StartDate:             entity.StartDate,
			EndDate:               entity.EndDate,
			TotalTodolist:         entity.TotalTodolist,
			TotalTodolistItemDone: entity.TotalTodolistItemDone,
			TotalTodolistItem:     entity.TotalTodolistItem,
			DaysRemaining:         daysRemaining,
			DaysRemainingStatus:   daysRemainingStatus,
			CompletionPercetage:   entity.CompletionPercentage,
		}
	}

	return dtos
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

	// Format tanggal dengan aman
	var startDate, endDate string

	if project.StartDate != nil {
		startDate = project.StartDate.Format("2006-01-02")
	}

	if project.EndDate != nil {
		endDate = project.EndDate.Format("2006-01-02")
	}

	return &dto.ProjectWithTodolistAndExpensesDto{
		ID:                         project.ID,
		UserID:                     project.UserID,
		Name:                       project.Name,
		Budget:                     project.Budget,
		StartDate:                  startDate,
		EndDate:                    endDate,
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

	todoIsAvailable := totalTodo > 0
	isTodolistComplete := todoIsAvailable && (totalTodo == totalCompletedTodo)
	return &dto.ProjectTodolistsDto{
		ID:                   todolist.ID,
		ProjectID:            todolist.ProjectID,
		Name:                 todolist.Name,
		TotalTodo:            totalTodo,
		TotalCompletedTodo:   totalCompletedTodo,
		ProjectTodolistItems: todolistItem,
		IsTodolistCompleted:  isTodolistComplete,
	}
}

func ToProjectTodolistItemDTO(todoItem *entities.ProjectTodolistItems) *dto.ProjectTodolistItemsDto {
	return &dto.ProjectTodolistItemsDto{
		ID:                todoItem.ID,
		ProjectTodolistID: todoItem.ProjectTodolistID,
		Name:              todoItem.Name,
		CategoryName:      todoItem.CategoryName,
		IsCompleted:       todoItem.IsCompleted,
	}
}

func ToProjectCategorySummaryDTO(projectCategory []*entities.ProjectCategorySummary) []*dto.ProjectCategorySummaryDto {

	var result = make([]*dto.ProjectCategorySummaryDto, 0, len(projectCategory)+1)

	totalAll := 0
	for _, item := range projectCategory {
		totalAll += item.Total
		result = append(result, &dto.ProjectCategorySummaryDto{
			CategoryName: item.CategoryName,
			Total:        item.Total,
		})
	}

	// Tambahkan kategori "All"
	result = append([]*dto.ProjectCategorySummaryDto{
		{
			CategoryName: "All",
			Total:        totalAll,
		},
	}, result...)

	return result
}

func ToProjectSummaryDTO(summary *entities.ProjectUserSummary) *dto.ProjectSummaryDto {

	return &dto.ProjectSummaryDto{
		TotalProjects:          summary.TotalProject,
		TotalBudgetUsed:        summary.TotalBudget,
		TotalCompletedProjects: summary.TotalCompletedProject,
	}
}
