package helper

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var secretKey = []byte(os.Getenv("JWT_SECRECT_KEY"))
var UserId uuid.UUID

func GenerateToken(userId string) (string, error) {
	// Log ketika proses token dimulai
	logrus.WithField("userId", userId).Info("Generating token for user")

	// Klaim token
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(), // Token valid selama 24 jam
	}

	// Log klaim yang digunakan
	logrus.WithFields(logrus.Fields{
		"userId": claims["userId"],
		"exp":    claims["exp"],
	}).Info("Claims set for token")

	// Membuat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		// Log error jika terjadi kesalahan saat menandatangani token
		logrus.WithError(err).Error("Failed to sign the token")
		return "", err
	}

	// Log keberhasilan pembuatan token
	logrus.WithField("token", tokenString).Info("Token generated successfully")
	return tokenString, nil
}

func VerifyToken(c *fiber.Ctx) error {
	tokenHeaderRaw := c.Get("Authorization", "")

	if tokenHeaderRaw == "" {
		logrus.Warn("Authorization header not provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Authorization header not provided",
		})
	}

	// Memisahkan "Bearer" dan token
	parts := strings.Split(tokenHeaderRaw, " ")
	if len(parts) != 2 || parts[1] == "" {
		logrus.Warn("Token not provided or invalid format")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Token not provided or invalid format",
		})
	}

	tokenHeader := parts[1]
	claims := jwt.MapClaims{}

	// Log parsing token dimulai
	logrus.Info("Parsing token...")
	token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Log ketika terjadi kesalahan dalam parsing token
	if err != nil {
		logrus.WithError(err).Error("Failed to parse token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to parse token",
		})
	}

	// Log jika token tidak valid
	if !token.Valid {
		logrus.Warn("Token invalid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Token invalid",
		})
	}

	// Log saat memeriksa klaim "userId"
	logrus.Info("Extracting userId from token claims...")
	if val, ok := claims["userId"].(string); ok {
		uuidUserId, errParseUserId := uuid.Parse(val)

		// Log kesalahan parsing userId
		if errParseUserId != nil {
			logrus.WithError(errParseUserId).Error("Failed to parse userId from token claims")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": errParseUserId.Error(),
			})
		}

		UserId = uuidUserId
		logrus.WithField("userId", uuidUserId).Info("Successfully extracted userId from token claims")

	} else {
		// Log jika userId tidak ditemukan di token
		logrus.Warn("User id not found in token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "User id not found in token",
		})
	}

	// Log jika token valid dan proses berhasil
	logrus.Info("Token verified successfully")
	return c.Next()
}
