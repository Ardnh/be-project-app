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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	err := h.projectsExpensesService.CreateProjectsExpenses(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusCreated, "Project expenses created", nil)
}
