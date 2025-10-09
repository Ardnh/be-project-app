package validation_utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Helper untuk format error messages
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			errors[field] = field + " is required"
		case "email":
			errors[field] = "Invalid email format"
		case "min":
			errors[field] = field + " must be at least " + err.Param() + " characters"
		case "max":
			errors[field] = field + " must be at most " + err.Param() + " characters"
		case "alphanum":
			errors[field] = field + " must contain only letters and numbers"
		case "eqfield":
			errors[field] = field + " must match " + err.Param()
		case "strongpassword":
			errors[field] = "Password must contain uppercase, lowercase, and number"
		default:
			errors[field] = "Invalid " + field
		}
	}

	return errors
}

// Helper untuk service errors
func handleServiceError(c *fiber.Ctx, err error) error {
	switch err.Error() {
	case "email already exists":
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	case "username already exists":
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already taken",
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
}
