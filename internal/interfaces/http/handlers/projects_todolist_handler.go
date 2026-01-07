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
	log.Printf("[INFO] Create project todolist - Method: %s, Path: %s, IP: %s", c.Method(), c.Path(), c.IP())

	var req dto.CreateProjectTodolistsDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v, Body: %s", err, string(c.Body()))
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Request body parsed successfully: %+v", req)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	log.Printf("[DEBUG] Validation passed, calling service layer")

	err := h.service.CreateProjectsTodolist(c.Context(), &req)
	if err != nil {
		log.Printf("[ERROR] Failed to create project todolist: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project todolist created successfully - ProjectId: %s, Name: %s", req.ProjectID, req.Name)
	return http.SuccessResponse(c, fiber.StatusCreated, "Project todolist created", nil)
}

func (h *ProjectTodolistHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Update project todolist - Method: %s, Path: %s, ID: %s, IP: %s", c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Update request without ID parameter")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectTodolistsDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body for ID %s: %v, Body: %s", id, err, string(c.Body()))
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Request body parsed successfully for ID %s: %+v", id, req)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed for ID %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	log.Printf("[DEBUG] Validation passed for ID %s, calling service layer", id)

	err := h.service.UpdateProjectsTodolist(c.Context(), id, &req)
	if err != nil {
		log.Printf("[ERROR] Failed to update project todolist %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project todolist updated successfully - ID: %s, Name: %s", id, req.Name)
	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist updated", nil)
}

func (h *ProjectTodolistHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")

	log.Printf("[INFO] Delete project todolist - Method: %s, Path: %s, ID: %s, IP: %s", c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Delete request without ID parameter")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	log.Printf("[DEBUG] Calling service layer to delete project todolist with ID: %s", id)

	err := h.service.DeleteProjectsTodolist(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete project todolist %s: %v", id, err)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Project todolist deleted successfully - ID: %s", id)
	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist deleted", nil)
}
