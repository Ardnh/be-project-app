package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProjectsExpensesItemHandler struct {
	projectsService services.ProjectsExpensesItemService
	validator       *validator.Validate
}

// Constructor - return *NewProjectsExpensesItemHandler (concrete type)
func NewProjectsExpensesItemHandler(projectsService services.ProjectsExpensesItemService, validator *validator.Validate) *ProjectsExpensesItemHandler {
	return &ProjectsExpensesItemHandler{
		projectsService: projectsService,
		validator:       validator,
	}
}

func (h *ProjectsExpensesItemHandler) Create(c *fiber.Ctx) error {

	var req dto.CreateProjectExpensesItemDto
	if err := c.BodyParser(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), "Invalid request body")
	}

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), validation_utils.FormatValidationErrors(err))
	}

	err := h.projectsService.CreateProjectsExpensesItem(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses item created", nil)
}

func (h *ProjectsExpensesItemHandler) Update(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectExpensesItemDto
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

	err := h.projectsService.UpdateProjectsExpensesItem(c.Context(), id, &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses item updated", nil)
}

func (h *ProjectsExpensesItemHandler) Delete(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	err := h.projectsService.DeleteProjectsExpensesItem(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses item deleted", nil)
}
