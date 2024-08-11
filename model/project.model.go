package model

import (
	"github.com/google/uuid"
)

type Projects struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Budget      int
}

type ProjectItem struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	Name       string
	BudgetItem int
	Status     bool
}

type CreateProjectRequest struct {
	CategoryID  uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Budget      int
}

type UpdateProjectRequest struct {
	CategoryID  uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Budget      int
}

type CreateProjectItem struct {
	ProjectID  uuid.UUID
	Name       string
	BudgetItem int
	Status     bool
}

type UpdateProjectItem struct {
	ProjectID  uuid.UUID
	Name       string
	BudgetItem int
	Status     bool
}

type ProjectItemResponse struct {
	Items   []ProjectItem
	Project Projects
}
