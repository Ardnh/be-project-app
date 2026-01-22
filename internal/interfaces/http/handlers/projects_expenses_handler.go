package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProjectsExpensesHandler struct {
	projectsExpensesService services.ProjectsExpensesService
	validator               *validator.Validate
}

// Constructor - return *UserHandler (concrete type)
func NewProjectsExpensesHandler(projectsExpensesService services.ProjectsExpensesService, validator *validator.Validate) *ProjectsExpensesHandler {
	return &ProjectsExpensesHandler{
		projectsExpensesService: projectsExpensesService,
		validator:               validator,
	}
}

func (h *ProjectsExpensesHandler) Create(c *fiber.Ctx) error {

	var req dto.CreateProjectsExpensesDto
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

	createdExpenseDto, err := h.projectsExpensesService.CreateProjectsExpenses(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusCreated, "Project expenses created", createdExpenseDto)
}

func (h *ProjectsExpensesHandler) Update(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectsExpensesDto
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

	updatedExpenseDto, err := h.projectsExpensesService.UpdateProjectsExpenses(c.Context(), id, &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses updated", updatedExpenseDto)
}

func (h *ProjectsExpensesHandler) Delete(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	err := h.projectsExpensesService.DeleteProjectsExpenses(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses deleted", nil)
}
