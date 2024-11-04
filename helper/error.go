package helper

import (
	"fmt"
	"project-app/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleValidationError[T any](err error, startTime time.Time, req T) fiber.Map {
	errorsField := []model.ErrorsField{}
	for _, err := range err.(validator.ValidationErrors) {
		logrus.WithFields(logrus.Fields{
			"field": err.Field(),
			"tag":   err.Tag(),
			"value": err.Value(),
		}).Warn("Field validation failed")

		errorsField = append(errorsField, model.ErrorsField{
			Field: err.Field(),
			Error: fmt.Sprintf("Field '%s' is not valid: %s", err.Field(), err.Tag()),
		})
	}

	return fiber.Map{
		"code":        fiber.StatusBadRequest,
		"message":     "Validation error",
		"processTime": ElapsedTime(startTime),
		"errors":      errorsField,
	}
}
