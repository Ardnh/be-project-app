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
	log.Printf("[INFO] Create todolist item - Method: %s, Path: %s, IP: %s", c.Method(), c.Path(), c.IP())

	var req dto.CreateProjectTodolistItemsDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse body: %v, Body: %s", err, string(c.Body()))
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), "Invalid request body")
	}

	log.Printf("[DEBUG] Request parsed: %+v", req)

	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), validation_utils.FormatValidationErrors(err))
	}

	log.Printf("[DEBUG] Validation passed, creating item")

	err := h.service.CreateProjectsTodolistItem(c.Context(), &req)
	if err != nil {
		log.Printf("[ERROR] Failed to create item: %v, Request: %+v", err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Todolist item created successfully - Name: %s, ProjectTodolistId: %s", req.Name, req.ProjectTodolistID)
	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist item created", nil)
}

func (h *ProjectTodolistItemHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id", "")
	log.Printf("[INFO] Update todolist item - Method: %s, Path: %s, ID: %s, IP: %s", c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Update request without ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	var req dto.UpdateProjectTodolistItemsDto
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse body for ID %s: %v, Body: %s", id, err, string(c.Body()))
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("[DEBUG] Request parsed for ID %s: %+v", id, req)

	if err := h.validator.Struct(&req); err != nil {
		log.Printf("[WARN] Validation failed for ID %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, err.Error(), fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	log.Printf("[DEBUG] Validation passed for ID %s, updating item", id)

	err := h.service.UpdateProjectsTodolistItem(c.Context(), id, &req)
	if err != nil {
		log.Printf("[ERROR] Failed to update item %s: %v, Request: %+v", id, err, req)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Todolist item updated successfully - ID: %s, Name: %s, IsCompleted: %v", id, req.Name, req.IsCompleted)
	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist item updated", nil)
}

func (h *ProjectTodolistItemHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	log.Printf("[INFO] Delete todolist item - Method: %s, Path: %s, ID: %s, IP: %s", c.Method(), c.Path(), id, c.IP())

	if id == "" {
		log.Printf("[WARN] Delete request without ID")
		return http.ErrorResponse(c, fiber.ErrBadRequest.Code, "ID not provided", nil)
	}

	log.Printf("[DEBUG] Deleting todolist item with ID: %s", id)

	err := h.service.DeleteProjectsTodolistItem(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete item %s: %v", id, err)
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	log.Printf("[INFO] Todolist item deleted successfully - ID: %s", id)
	return http.SuccessResponse(c, fiber.StatusOK, "Project todolist item deleted", nil)
}
