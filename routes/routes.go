package routes

import (
	"fmt"
	"project-app/handler/category"
	"project-app/handler/users"
	"project-app/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, validate *validator.Validate) {

	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("Request: %s %s \n", c.Method(), c.OriginalURL())
		return c.Next()
	})

	userHandler := users.NewUsersHandler(db, validate)
	categoryHandler := category.NewCategoryHandler(db, validate)

	appGroup := app.Group("/api/v1")

	// Swagger
	appGroup.Get("/swagger/*", swagger.HandlerDefault)

	// Users
	usersGroup := appGroup.Group("user")
	usersGroup.Post("/login", userHandler.Login)
	usersGroup.Post("/register", userHandler.Register)
	usersGroup.Get("/profile/:user_id", helper.VerifyToken, userHandler.GetProfileById)
	usersGroup.Put("/profile", helper.VerifyToken, userHandler.UpdateProfileById)

	// Category
	categoryGroup := appGroup.Group("category")
	categoryGroup.Post("/", categoryHandler.Create)
	categoryGroup.Put("/", categoryHandler.Update)
	categoryGroup.Delete("/:id", categoryHandler.Delete)
	categoryGroup.Get("/", categoryHandler.FindAll)

}
