package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func (h *ProjectTodolistHandler) Create(c *fiber.Ctx) error {

	var req dto.CreateProjectTodolistsDto
	if err := c.BodyParser(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	err := h.service.CreateProjectsTodolist(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusCreated, "Project todolist created", nil)
}

func (h *ProjectTodolistHandler) Update(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectTodolistsDto
	if err := c.BodyParser(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	err := h.service.UpdateProjectsTodolist(c.Context(), id, &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist updated", nil)
}

func (h *ProjectTodolistHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	err := h.service.DeleteProjectsTodolist(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses deleted", nil)
}
