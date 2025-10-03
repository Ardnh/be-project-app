package handlers

import (
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
)

// UserHandler struct - NO INTERFACE NEEDED (opsional)
type UserHandler struct {
	userService services.UserService // ✅ Interface dari service layer
	validator   *validator.Validate
}

// Constructor - return *UserHandler (concrete type)
func NewUserHandler(userService services.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {

	return nil
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{"data": "Hello"})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {

	return nil
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {

	return nil
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {

	return nil
}
