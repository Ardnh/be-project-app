package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProjectsHandler struct {
	projectsService services.ProjectsService // ✅ Interface dari service layer
	validator       *validator.Validate
}

// Constructor - return *UserHandler (concrete type)
func NewProjectsHandler(projectsService services.ProjectsService, validator *validator.Validate) *ProjectsHandler {
	return &ProjectsHandler{
		projectsService: projectsService,
		validator:       validator,
	}
}

func (h *ProjectsHandler) CreateProject(c *fiber.Ctx) error {
	var req dto.CreateProjectsDto
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

	err := h.projectsService.CreateProjects(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusCreated, "User created", nil)
}

func (h *ProjectsHandler) GetProjectsByUserId(c *fiber.Ctx) error {

	user_id := c.Params("user_id", "")
	if user_id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "User ID not provided", nil)
	}

	projects, err := h.projectsService.GetProjectsByUserID(c.Context(), user_id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Success Get Projects", projects)
}
