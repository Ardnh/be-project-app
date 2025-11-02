// internal/interfaces/http/handlers/health_handler.go
package handlers

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(dto.Success("Service is healthy", fiber.Map{
		"status":  "UP",
		"service": "project-app",
	}))
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Check if service is ready to accept requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *fiber.Ctx) error {
	// TODO: Check database, redis, etc.
	return c.JSON(dto.Success("Service is ready", fiber.Map{
		"database": "connected",
		"redis":    "connected",
	}))
}
