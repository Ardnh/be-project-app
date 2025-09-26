package http

import (
	"errors"
	"net/http"

	"github.com/Ardnh/be-project-app/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success response
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Handle service layer errors
func HandleServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return ErrorResponse(c, http.StatusNotFound, "User not found", nil)
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return ErrorResponse(c, http.StatusConflict, "Email already exists", nil)
	case errors.Is(err, domain.ErrUnauthorized):
		return ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
	default:
		return ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
	}
}
