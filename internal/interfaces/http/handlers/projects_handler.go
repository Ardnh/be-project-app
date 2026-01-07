package handlers

import (
	"errors"
	"log"
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
	log.Printf("[INFO] Create project - Method: %s, Path: %s, IP: %s", c.Method(), c.Path(), c.IP())

	var req dto.CreateProjectDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v, Body: %s", err, string(c.Body()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Request body parsed successfully: %+v", req)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed: %v, Request: %+v", err, req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	log.Printf("[DEBUG] Validation passed, calling service layer")

	err := h.projectsService.CreateProjects(c.Context(), &req)
	if err != nil {
		log.Printf("[ERROR] Failed to create project: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project created successfully - Name: %s, UserId: %s", req.Name, req.UserID)
	return http.SuccessResponse(c, fiber.StatusCreated, "Success create project", nil)
}

func (h *ProjectsHandler) GetProjectsByUserId(c *fiber.Ctx) error {
	userId := c.Params("user_id", "")

	log.Printf("[INFO] Get projects by user ID - Method: %s, Path: %s, UserID: %s, IP: %s",
		c.Method(), c.Path(), userId, c.IP())

	if userId == "" {
		log.Printf("[WARN] Get projects request without user ID")
		return http.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	// Parse query parameters
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var offset int
	pageStr := c.Query("page", "1")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		offset = (page - 1) * limit
	}

	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "DESC")

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	projectName := c.Query("search", "")

	params := dto.GetProjectsParams{
		UserID:    userId,
		Search:    projectName,
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	log.Printf("[DEBUG] Query params - UserID: %s, Search: %s, Limit: %d, Offset: %d, SortBy: %s, SortOrder: %s",
		userId, projectName, limit, offset, sortBy, sortOrder)

	result, err := h.projectsService.GetProjectsByUserID(c.Context(), userId, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[WARN] User not found: %s", userId)
			return http.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
		}
		log.Printf("[ERROR] Failed to get projects for user %s: %v", userId, err)
		return http.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get projects", err.Error())
	}

	log.Printf("[INFO] Successfully retrieved projects for user %s - Count: %d", userId, len(result))
	return http.SuccessResponse(c, fiber.StatusOK, "Success get projects", result)
}

func (h *ProjectsHandler) GetProjectById(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Get project by ID - Method: %s, Path: %s, ID: %s, IP: %s",
		c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Get project request without ID")
		return http.ErrorResponse(c, fiber.StatusBadRequest, "Project ID is required", nil)
	}

	log.Printf("[DEBUG] Fetching project with ID: %s", id)

	project, err := h.projectsService.GetProjectsByID(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Failed to get project %s: %v", id, err)
		return http.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get project", err.Error())
	}

	log.Printf("[INFO] Successfully retrieved project - ID: %s, Name: %s", id, project.Name)
	return http.SuccessResponse(c, fiber.StatusOK, "Success get project", project)
}

func (h *ProjectsHandler) GetAllProjectSummaryByUserId(c *fiber.Ctx) error {
	userId := c.Params("user_id", "")

	log.Printf("[INFO] Get project summary by user ID - Method: %s, Path: %s, UserID: %s, IP: %s",
		c.Method(), c.Path(), userId, c.IP())

	if userId == "" {
		log.Printf("[WARN] Get summary request without user ID")
		return http.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	log.Printf("[DEBUG] Parsing user ID: %s", userId)

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		log.Printf("[ERROR] Invalid user ID format: %s, Error: %v", userId, err)
		return http.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}

	log.Printf("[DEBUG] Fetching summary for user: %s", parsedUserId.String())

	summary, err := h.projectsService.GetProjectsSummaryByUserID(c.Context(), (parsedUserId).String())
	if err != nil {
		log.Printf("[ERROR] Failed to get summary for user %s: %v", userId, err)
		return http.ErrorResponse(c, fiber.StatusBadRequest, "Failed to get summary", err)
	}

	log.Printf("[INFO] Successfully retrieved summary for user %s", userId)
	return http.SuccessResponse(c, fiber.StatusOK, "Successfully get summary", summary)
}

func (h *ProjectsHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Update project - Method: %s, Path: %s, ID: %s, IP: %s",
		c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Update project request without ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectDto

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body for project %s: %v, Body: %s", id, err, string(c.Body()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Request body parsed successfully for project %s: %+v", id, req)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed for project %s: %v, Request: %+v", id, err, req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	log.Printf("[DEBUG] Validation passed for project %s, calling service layer", id)

	err := h.projectsService.UpdateProjects(c.Context(), id, &req)
	if err != nil {
		log.Printf("[ERROR] Failed to update project %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project updated successfully - ID: %s, Name: %s", id, req.Name)
	return http.SuccessResponse(c, fiber.StatusOK, "Successfully update project", nil)
}

func (h *ProjectsHandler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Delete project - Method: %s, Path: %s, ID: %s, IP: %s",
		c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Delete project request without ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	log.Printf("[DEBUG] Calling service layer to delete project with ID: %s", id)

	err := h.projectsService.DeleteProjects(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete project %s: %v", id, err)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project deleted successfully - ID: %s", id)
	return http.SuccessResponse(c, fiber.StatusOK, "Successfully deleted project", nil)
}

func (s *ProjectsHandler) GetProjectCategoryByUserId(c *fiber.Ctx) error {
	id := c.Params("user_id", "")

	log.Printf("[INFO] Get project category by user ID - Method: %s, Path: %s, UserID: %s, IP: %s",
		c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Get category request without user ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	log.Printf("[DEBUG] Fetching project categories for user: %s", id)

	result, err := s.projectsService.GetProjectCategoryByUserId(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Failed to get project category for user %s: %v", id, err)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Successfully retrieved project categories for user %s - Count: %d", id, len(result))
	return http.SuccessResponse(c, fiber.StatusOK, "Successfully get project category", result)
}
