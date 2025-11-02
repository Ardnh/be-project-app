package handlers

import (
	"strings"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

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

// GetUserByToken handles user profile retrieval from JWT token
func (h *UserHandler) GetUserByToken(c *fiber.Ctx) error {
	// 1. Ambil JWT claims dari context
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
	}

	// 2. Ambil user_id dari claims
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token payload", nil)
	}

	// 3. Ambil user dari service layer
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return http.HandleServiceError(c, err)
	}

	// 5. Return success response
	return http.SuccessResponse(c, fiber.StatusOK, "Successfully retrieved user", user)
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
