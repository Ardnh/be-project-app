package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// TotalProcessTime adalah middleware yang menghitung waktu proses setiap request
func ProcessTime() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Proses request selanjutnya dalam chain middleware/handler
		err := c.Next()

		// Hitung durasi setelah request selesai diproses
		duration := time.Since(start)

		// Log informasi eksekusi
		logrus.WithFields(logrus.Fields{
			"method":        c.Method(),
			"route":         c.Path(),
			"status":        c.Response().StatusCode(),
			"executionTime": duration.Milliseconds(), // Menggunakan Milliseconds
		}).Info("Request processed")

		// Simpan durasi dalam c.Locals agar bisa diakses di handler jika diperlukan
		c.Locals("executionTime", duration.Milliseconds())

		return err
	}
}
