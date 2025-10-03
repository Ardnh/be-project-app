// internal/interfaces/http/middlewares/auth.go
package middlewares

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		token := c.Get("Authorization")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				dto.Error("Unauthorized: No token provided"),
			)
		}

		// TODO: Validate JWT token
		// user, err := jwt.ValidateToken(token)
		// if err != nil {
		//     return c.Status(fiber.StatusUnauthorized).JSON(...)
		// }

		// Store user in context
		// c.Locals("user", user)

		return c.Next()
	}
}

// AdminMiddleware validates if user is admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Check if user is admin
		// user := c.Locals("user")
		// if !user.IsAdmin {
		//     return c.Status(fiber.StatusForbidden).JSON(...)
		// }

		return c.Next()
	}
}
