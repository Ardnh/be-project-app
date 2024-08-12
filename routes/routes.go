package routes

import (
	"fmt"
	"project-app/handler/category"
	proeject "project-app/handler/project"
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
	projectHandler := proeject.NewProjectHandler(db, validate)

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

	// Project
	projectGroup := appGroup.Group("projects")
	projectGroup.Post("/", projectHandler.CreateProject)
	projectGroup.Put("/:id", projectHandler.UpdateProject)
	projectGroup.Delete("/:id", projectHandler.DeleteProject)
	projectGroup.Get("/", projectHandler.GetAllProject)

	// Project item
	projectItemGroup := appGroup.Group("project-item")
	projectItemGroup.Post("/", projectHandler.CreateProjectItem)
	projectGroup.Put("/:id", projectHandler.UpdateProjectItem)
	projectGroup.Delete("/:id", projectHandler.DeleteProjectItem)
	projectGroup.Get("/:project_id", projectHandler.GetAllProjectItemByProjectId)
}
