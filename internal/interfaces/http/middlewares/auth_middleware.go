// internal/interfaces/http/middlewares/auth.go
package middlewares

import (
	"strings"
	"time"

	"github.com/Ardnh/be-project-app/internal/config"
	http "github.com/Ardnh/be-project-app/internal/interfaces/http/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT token and attaches user info to context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return http.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: No token provided", nil)
		}

		// Expect header format: "Bearer <token>"
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			return http.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token format", nil)
		}

		// Load JWT secret key from config
		cfg := config.LoadConfig()
		secretKey := []byte(cfg.App.JWTSecret)
		if len(secretKey) == 0 {
			return http.ErrorResponse(c, fiber.StatusInternalServerError, "JWT secret not configured", nil)
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			// Ensure token uses correct signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
		}

		// Validate claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Optional: check expiry manually
			if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
				return http.ErrorResponse(c, fiber.StatusUnauthorized, "Token has expired", nil)
			}

			// Store user info in context
			c.Locals("user", claims)
			return c.Next()
		}

		return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
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
