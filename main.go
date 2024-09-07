package main

import (
	"os"
	"project-app/app"
	"project-app/helper"
	"project-app/routes"
	"project-app/config"

	_ "project-app/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

// @title           Project APP API
// @version         1.0
// @description     API Documentation for Project APP API.

// @contact.name   Muhammad Ardan Hilal
// @contact.url    ardn.h79@gmail.com
// @contact.email  ardn.h79@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {

	newApp := fiber.New()
	db := app.DbConnection()
	validate := validator.New(validator.WithRequiredStructEnabled())

	config.LoggerConfig()
	routes.SetupRoutes(newApp, db, validate)

	newApp.Use(config.LogrusMiddleware)

	port := os.Getenv("APP_PORT")
	err := newApp.Listen(port)
	helper.PanicIfError(err)
}
