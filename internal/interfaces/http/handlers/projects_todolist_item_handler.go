package handlers

import (
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"github.com/go-playground/validator/v10"
)

type ProjectTodolistItemHandler struct {
	service   services.ProjectTodolistItemService
	validator *validator.Validate
}

func NewProjectTodolistItemHandler(service services.ProjectTodolistItemService, validator *validator.Validate) *ProjectTodolistItemHandler {
	return &ProjectTodolistItemHandler{
		service:   service,
		validator: validator,
	}
}
