package handlers

import (
	"log"

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
	log.Printf("[INFO] Create project expenses item - Method: %s, Path: %s, IP: %s", c.Method(), c.Path(), c.IP())

	var req dto.CreateProjectExpensesItemDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v, Body: %s", err, string(c.Body()))
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), "Invalid request body")
	}

	log.Printf("[DEBUG] Request body parsed successfully: %+v", req)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), validation_utils.FormatValidationErrors(err))
	}

	log.Printf("[DEBUG] Validation passed, calling service layer")

	err := h.projectsService.CreateProjectsExpensesItem(c.Context(), &req)
	if err != nil {
		log.Printf("[ERROR] Failed to create expenses item: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project expenses item created successfully - ProjectExpensesId: %s, Name: %s, Amount: %f",
		req.ProjectExpensesId, req.Name, req.Amount)
	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses item created", nil)
}

func (h *ProjectsExpensesItemHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Update project expenses item - Method: %s, Path: %s, ID: %s, IP: %s",
		c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Update expenses item request without ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectExpensesItemDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body for expenses item %s: %v, Body: %s", id, err, string(c.Body()))
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Request body parsed successfully for expenses item %s: %+v", id, req)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed for expenses item %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	log.Printf("[DEBUG] Validation passed for expenses item %s, calling service layer", id)

	err := h.projectsService.UpdateProjectsExpensesItem(c.Context(), id, &req)
	if err != nil {
		log.Printf("[ERROR] Failed to update expenses item %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project expenses item updated successfully - ID: %s, Name: %s, Amount: %f",
		id, req.Name, req.Amount)
	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses item updated", nil)
}

func (h *ProjectsExpensesItemHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Delete project expenses item - Method: %s, Path: %s, ID: %s, IP: %s",
		c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Delete expenses item request without ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	log.Printf("[DEBUG] Calling service layer to delete expenses item with ID: %s", id)

	err := h.projectsService.DeleteProjectsExpensesItem(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete expenses item %s: %v", id, err)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project expenses item deleted successfully - ID: %s", id)
	return http.SuccessResponse(c, fiber.StatusOK, "Project expenses item deleted", nil)
}
