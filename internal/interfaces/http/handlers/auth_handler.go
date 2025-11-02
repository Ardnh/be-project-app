package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	validation_utils "github.com/Ardnh/be-project-app/internal/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlers struct {
	authService services.AuthService // ✅ Interface dari service layer
	validator   *validator.Validate
}

func NewAuthHandlers(authService services.AuthService, validator *validator.Validate) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		validator:   validator,
	}
}

func (h *AuthHandlers) Login(c *fiber.Ctx) error {

	var req dto.LoginDto
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

	tokenString, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return http.HandleServiceError(c, err)
	}

	res := map[string]any{
		"token": tokenString,
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Login successfully", res)
}

func (h *AuthHandlers) Register(c *fiber.Ctx) error {

	var req dto.RegisterDto
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

	err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return http.ErrorResponse(c, fiber.ErrInternalServerError.Code, err.Error(), nil)
	}

	return http.SuccessResponse(c, fiber.StatusOK, "Register successfully", nil)
}
