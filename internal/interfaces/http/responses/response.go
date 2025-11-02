package http

import (
	"errors"
	"net/http"

	"github.com/Ardnh/be-project-app/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
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
		return ErrorResponse(c, http.StatusNotFound, domain.ErrInvalidCredentials.Error(), nil)
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return ErrorResponse(c, http.StatusConflict, domain.ErrEmailAlreadyExists.Error(), nil)
	case errors.Is(err, domain.ErrUnauthorized):
		return ErrorResponse(c, http.StatusUnauthorized, domain.ErrUnauthorized.Error(), nil)
	default:
		return ErrorResponse(c, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), nil)
	}
}
