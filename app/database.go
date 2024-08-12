package app

import (
	"log"
	"os"
	"project-app/helper"
	"project-app/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {

	postgresDsn := os.Getenv("APP_DSN")
	seedUsername := os.Getenv("APP_SEED_USERNAME")
	seedEmail := os.Getenv("APP_SEED_EMAIL")
	seedPassword := os.Getenv("APP_SEED_PASSWORD")
	db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})

	helper.PanicIfError(err)

	// Menjalankan perintah SQL untuk membuat ekstensi uuid-ossp
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatal("Failed to create uuid-ossp extension:", err)
	}

	db.AutoMigrate(
		&schema.Category{},
		&schema.User{},
		&schema.Project{},
		&schema.ProjectItem{},
	)

	// RUN SEEDER
	hashedPassword, errHash := helper.HashPassword(seedPassword)

	if errHash != nil {
		log.Fatal("Failed to hash password: ", err)
	}

	if errSeeder := db.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?) ON CONFLICT (email) DO NOTHING`, seedUsername, seedEmail, hashedPassword).Error; errSeeder != nil {
		log.Fatal("Failed to insert seed user: ", errSeeder)
	}

	return db
}
