package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		log.WithFields(logrus.Fields{
			"type":     "http",
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   c.Response().StatusCode(), // int
			"duration": duration.String(),
			"ip":       c.IP(),
		}).Info("request")

		return err
	}
}
