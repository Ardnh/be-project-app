package services

import (
	repositories "github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/gofiber/fiber/v2"
)

type userService struct {
	userRepo repositories.UserRepository
}

func GetUsers(c *fiber.Ctx) {

}

func CreateUser(c *fiber.Ctx) {

}

func UpdateUser(c *fiber.Ctx) {

}

func DeleteUser(c *fiber.Ctx) {

}
