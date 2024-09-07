package config

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggerConfig() {
	// Setup Logrus
	log := logrus.New()
	log.SetOutput(os.Stdout)       // Output log ke terminal
	log.SetLevel(logrus.InfoLevel) // Set log level
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func LogrusMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process the request
	err := c.Next()

	// Logrus logging
	logrus.WithFields(logrus.Fields{
		"status":     c.Response().StatusCode(),
		"method":     c.Method(),
		"path":       c.Path(),
		"ip":         c.IP(),
		"latency":    time.Since(start).String(),
		"user-agent": c.Get("User-Agent"),
	}).Info("Handled request")

	return err
}
