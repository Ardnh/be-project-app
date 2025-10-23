package dto

// ================ DTO ====================
type ProjectsDto struct {
	ID           string  `json:"id" validate:"required"`
	UserID       string  `json:"user_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
	CategoryName string  `json:"category_name" validate:"required"`
}

type ProjectWithTodolistAndExpensesDto struct {
	ID              string                 `json:"id" validate:"required"`
	UserID          string                 `json:"user_id" validate:"required"`
	Name            string                 `json:"name" validate:"required"`
	Budget          float64                `json:"budget" validate:"required"`
	StartDate       string                 `json:"start_date"`
	EndDate         string                 `json:"end_date"`
	CategoryName    string                 `json:"category_name" validate:"required"`
	ProjectExpenses []*ProjectsExpensesDto `json:"project_expenses" validate:"required"`
	ProjectTodolist []*ProjectTodolistsDto `json:"project_todolists" validate:"required"`
}

type ProjectCategorySummaryDto struct {
	CategoryName string `json:"category_name" validate:"required"`
	Total        int    `json:"total" validate:"required"`
}

type ProjectSummaryDto struct {
	TotalBudgetUsed   int `json:"total_budget_used" validate:"required"`
	TotalProjects     int `json:"total_projects" validate:"required"`
	TotalProjectsDone int `json:"total_projects_done" validate:"required"`
}

// ================ REQUEST DTO ====================
type CreateProjectsDto struct {
	UserID       string  `json:"user_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	CategoryName string  `json:"category_name" validate:"required"`
}

type UpdateProjectsDto struct {
	UserID       string  `json:"user_id" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	CategoryName string  `json:"category_name" validate:"required"`
}

// ================ RESPONSE DTO ====================
