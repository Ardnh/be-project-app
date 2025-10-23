package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func (h *ProjectTodolistItemHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProjectTodolistItemsDto
	if err := c.BodyParser(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), "Invalid request body")
	}

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), validation_utils.FormatValidationErrors(err))
	}

	err := h.service.CreateProjectsTodolistItem(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist item created", nil)
}

func (h *ProjectTodolistItemHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectTodolistItemsDto
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

	err := h.service.UpdateProjectsTodolistItem(c.Context(), id, &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist item updated", nil)
}

func (h *ProjectTodolistItemHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	err := h.service.DeleteProjectsTodolistItem(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist item deleted", nil)
}
