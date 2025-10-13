// cmd/api/main.go
package main

import (
	"fmt"
	"log"

	"github.com/Ardnh/be-project-app/internal/application/services"
	"github.com/Ardnh/be-project-app/internal/infrastructure/cache/redis"
	"github.com/Ardnh/be-project-app/internal/infrastructure/database/postgresql"
	"github.com/Ardnh/be-project-app/internal/infrastructure/database/repository"
	"github.com/Ardnh/be-project-app/internal/interfaces/http/handlers"
	"github.com/Ardnh/be-project-app/internal/interfaces/http/routes"
	"github.com/go-playground/validator/v10"

	"github.com/Ardnh/be-project-app/internal/config"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db, err := postgresql.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer postgresql.CloseDB(db)

	redisDb := redis.NewRedisDB(cfg)
	defer redisDb.Close()

	// TODO: Setup routes, handlers, etc.
	app := fiber.New()
	validator := validator.New()

	// Repository | interface -> infrastructure -> database -> repository
	userRepository := repository.NewUserRepository(db, redisDb)
	projectsRepository := repository.NewProjectsRepository(db, redisDb)
	projectsExpensesRepository := repository.NewProjectsExpensesRepository(db, redisDb)

	// Service | internal -> application -> service
	userService := services.NewUserService(userRepository)
	projectsService := services.NewProjectsService(projectsRepository)
	projectsExpensesService := services.NewProjectExpensesService(projectsExpensesRepository)

	// Handler | internal -> interfaces -> http -> handler
	userHandler := handlers.NewUserHandler(userService, validator)
	projectsHandler := handlers.NewProjectsHandler(projectsService, validator)
	projectsExpensesHandler := handlers.NewProjectsExpensesHandler(projectsExpensesService, validator)
	healthHandler := handlers.NewHealthHandler()

	// Setup Routes
	routes.SetupAPIRoutes(
		app,
		userHandler,
		projectsHandler,
		projectsExpensesHandler,
	)
	routes.SetupHealthRoutes(app, healthHandler)

	log.Printf("🚀 Starting application in %s mode on port %s\n", cfg.App.Env, cfg.App.Port)
	log.Println("✅ Application started successfully!")

	portListen := fmt.Sprintf(":%s", cfg.App.Port)
	if err := app.Listen(portListen); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
