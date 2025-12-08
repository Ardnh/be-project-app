// internal/interfaces/http/routes/health.go
package routes

import (
	"github.com/Ardnh/be-project-app/internal/interfaces/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupHealthRoutes(app *fiber.App, handler *handlers.HealthHandler) {
	app.Get("/health", handler.HealthCheck)
}
