package handlers

import (
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"github.com/go-playground/validator/v10"
)

type ProjectTodolistHandler struct {
	service   services.ProjectTodolistService
	validator *validator.Validate
}

func NewProjectTodolistHandler(service services.ProjectTodolistService, validator *validator.Validate) *ProjectTodolistHandler {
	return &ProjectTodolistHandler{
		service:   service,
		validator: validator,
	}
}
