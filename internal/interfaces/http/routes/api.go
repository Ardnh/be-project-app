// internal/interfaces/http/routes/api.go
package routes

import (
	"github.com/Ardnh/be-project-app/internal/interfaces/http/handlers"
	"github.com/Ardnh/be-project-app/internal/interfaces/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupAPIRoutes(
	app *fiber.App,
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectsHandler,
	projectExpensesHandler *handlers.ProjectsExpensesHandler,
	projectExpensesItemHandler *handlers.ProjectsExpensesItemHandler,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Public routes
	// public := api.Group("/")
	{
		// Auth routes (nanti)
		// public.Post("/login", authHandler.Login)
		// public.Post("/register", authHandler.Register)
	}

	// Protected routes (require authentication)
	protected := api.Group("/", middlewares.AuthMiddleware())
	{
		// User routes
		users := protected.Group("/users")
		{
			users.Get("/", userHandler.GetAllUsers)
			users.Get("/:id", userHandler.GetUserByID)
			users.Post("/", userHandler.CreateUser)
			users.Put("/:id", userHandler.UpdateUser)
			users.Delete("/:id", userHandler.DeleteUser)
		}

		// Project routes (contoh)
		projects := protected.Group("/projects")
		{
			projects.Get("/:user_id", projectHandler.GetProjectsByUserId)
			projects.Get("/by-id/:id", projectHandler.GetProjectById)
			projects.Get("/category/:user_id", projectHandler.GetProjectCategoryByUserId)
			projects.Get("/summary/:user_id", projectHandler.GetAllProjectSummaryByUserId)
			projects.Post("/", projectHandler.CreateProject)
			projects.Put("/:id", projectHandler.UpdateProject)
			projects.Delete("/:id", projectHandler.DeleteProject)
			// projects.Get("/:id", projectHandler.GetProjectByID)
		}

		// Project Expenses
		projectsExpenses := protected.Group("/project-expenses")
		{
			projectsExpenses.Post("/", projectExpensesHandler.Create)
			projectsExpenses.Put("/:id", projectExpensesHandler.Update)
			projectsExpenses.Delete("/:id", projectExpensesHandler.Delete)
		}

		// Project Expenses Item
		projectExpensesItem := protected.Group("/project-expenses-item")
		{
			projectExpensesItem.Post("/", projectExpensesHandler.Create)
			projectExpensesItem.Put("/:id", projectExpensesHandler.Update)
			projectExpensesItem.Delete("/:id", projectExpensesHandler.Delete)
		}
	}

	// Admin routes (require admin role)
	admin := api.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.Get("/users", userHandler.GetAllUsers)
		// admin.Delete("/users/:id", userHandler.DeleteUser)
	}
}
