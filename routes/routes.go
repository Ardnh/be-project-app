package routes

import (
	"project-app/handler/category"
	project "project-app/handler/project"
	"project-app/handler/users"
	"project-app/helper"
	categoryRepository "project-app/repository/category"
	projectRepository "project-app/repository/project"
	userRepository "project-app/repository/users"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, validate *validator.Validate) {

	app.Use(func(c *fiber.Ctx) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Method(),
			"url":    c.OriginalURL(),
			"ip":     c.IP(),
		}).Info("Incoming request")
		return c.Next()
	})

	// Repository
	userRepository := userRepository.NewUsersRepository(db)
	projectRepository := projectRepository.NewProjectRepository(db)
	categoryRepository := categoryRepository.NewCategoryRepository(db)

	// Handler
	userHandler := users.NewUsersHandler(userRepository, validate)
	categoryHandler := category.NewCategoryHandler(categoryRepository, validate)
	projectHandler := project.NewProjectHandler(projectRepository, validate)

	// Route
	appGroup := app.Group("/api/v1")

	// Swagger
	appGroup.Get("/swagger/*", swagger.HandlerDefault)

	// Users
	usersGroup := appGroup.Group("/user")
	usersGroup.Post("/login", userHandler.Login)
	usersGroup.Post("/register", userHandler.Register)
	usersGroup.Get("/profile/:user_id", helper.VerifyToken, userHandler.GetProfileById)
	usersGroup.Put("/profile", helper.VerifyToken, userHandler.UpdateProfileById)

	// Category
	categoryGroup := appGroup.Group("/category")
	categoryGroup.Post("/", helper.VerifyToken, categoryHandler.Create)
	categoryGroup.Put("/:id", helper.VerifyToken, categoryHandler.Update)
	categoryGroup.Delete("/:id", helper.VerifyToken, categoryHandler.Delete)
	categoryGroup.Get("/", categoryHandler.FindAll)

	// Project
	projectGroup := appGroup.Group("/projects")
	projectGroup.Post("/", projectHandler.CreateProject)
	projectGroup.Put("/:id", projectHandler.UpdateProject)
	projectGroup.Delete("/:id", projectHandler.DeleteProject)
	projectGroup.Get("/", projectHandler.GetAllProject)

	// Project item
	projectItemGroup := appGroup.Group("/project-item")
	projectItemGroup.Post("/", projectHandler.CreateProjectItem)
	projectItemGroup.Put("/:id", projectHandler.UpdateProjectItem)
	projectItemGroup.Delete("/:id", projectHandler.DeleteProjectItem)
	projectItemGroup.Get("/:project_id", projectHandler.GetAllProjectItemByProjectId)
}
