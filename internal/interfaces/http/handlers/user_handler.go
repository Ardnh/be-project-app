package handlers

import (
	"strings"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
)

// UserHandler struct - NO INTERFACE NEEDED (opsional)
type UserHandler struct {
	userService services.UserService // ✅ Interface dari service layer
	validator   *validator.Validate
}

// Constructor - return *UserHandler (concrete type)
func NewUserHandler(userService services.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// 1. Ambil ID dari URL params
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// 2. Parse UUID (jika menggunakan UUID)
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// 3. Call service
	user, err := h.userService.GetUserByID(c.Context(), userID.String())
	if err != nil {
		// Handle error based on type
		if err.Error() == "user not found" { // atau gunakan custom error
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	// 4. Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{"data": "Hello"})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Trim spaces
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	// Validate
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validation_utils.FormatValidationErrors(err),
		})
	}

	err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusCreated, "User created", nil)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {

	return nil
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {

	return nil
}
