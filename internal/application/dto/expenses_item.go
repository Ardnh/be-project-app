package dto

// ================ DTO ====================
type ProjectExpensesItemDto struct {
	ID                string  `json:"id" validate:"required"`
	ProjectExpensesId string  `json:"project_expenses_id" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
	CategoryName      string  `json:"category_name" validate:"required"`
}

// ================ REQUEST DTO ====================
type CreateProjectExpensesItemDto struct {
	ProjectExpensesId string  `json:"project_expense_id" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
	CategoryName      string  `json:"category_name" validate:"required"`
}

type UpdateProjectExpensesItemDto struct {
	ProjectExpensesId string  `json:"project_expense_id" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
	CategoryName      string  `json:"category_name" validate:"required"`
}
