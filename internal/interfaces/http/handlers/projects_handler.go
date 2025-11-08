package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectsHandler struct {
	projectsService services.ProjectsService
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
	var req dto.CreateProjectDto
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

	return http.SuccessResponse(c, fiber.StatusCreated, "Success create project", nil)
}

func (h *ProjectsHandler) GetProjectsByUserId(c *fiber.Ctx) error {
	// 1. Get user_id dari URL params
	userId := c.Params("user_id", "")
	if userId == "" {
		return http.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	// 2. Parse query parameters dengan default values
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50 // Max limit untuk prevent abuse
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "DESC")

	// Validate sortOrder
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Validate sortBy (whitelist allowed fields)
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	// 3. Optional: Get additional filters
	projectName := c.Query("search", "")

	// 4. Build request DTO/params
	params := dto.GetProjectsParams{
		UserID:    userId,
		Search:    projectName,
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// 5. Call service
	result, err := h.projectsService.GetProjectsByUserID(c.Context(), userId, params)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
		}
		return http.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get projects", err.Error())
	}

	// 6. Return success response
	return http.SuccessResponse(c, fiber.StatusOK, "Success get projects", result)
}

func (h *ProjectsHandler) GetProjectById(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.StatusBadRequest, "Project ID is required", nil)
	}

	project, err := h.projectsService.GetProjectsByID(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get project", err.Error())
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Success get project", project)
}

func (h *ProjectsHandler) GetAllProjectSummaryByUserId(c *fiber.Ctx) error {

	// 1. Get user_id dari URL params
	userId := c.Params("user_id", "")
	if userId == "" {
		return http.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return http.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}

	summary, err := h.projectsService.GetProjectsSummaryByUserID(c.Context(), (parsedUserId).String())
	if err != nil {
		return http.ErrorResponse(c, fiber.StatusBadRequest, "Failed to get summary", err)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Successfully get summary", summary)
}

func (h *ProjectsHandler) UpdateProject(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectDto

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

	err := h.projectsService.UpdateProjects(c.Context(), id, &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Successfully update project", nil)
}

func (h *ProjectsHandler) DeleteProject(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	err := h.projectsService.DeleteProjects(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Successfully deleted project", nil)
}

func (s *ProjectsHandler) GetProjectCategoryByUserId(c *fiber.Ctx) error {

	id := c.Params("user_id", "")
	if id == "" {
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	result, err := s.projectsService.GetProjectCategoryByUserId(c.Context(), id)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Successfully get project category", result)
}
