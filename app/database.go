package app

import (
	"github.com/sirupsen/logrus"
	"os"
	"project-app/helper"
	"project-app/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {
    // Ambil variabel environment
    postgresDsn := os.Getenv("APP_DSN")
    seedUsername := os.Getenv("APP_SEED_USERNAME")
    seedEmail := os.Getenv("APP_SEED_EMAIL")
    seedPassword := os.Getenv("APP_SEED_PASSWORD")

    logrus.Info("Connecting to database...")

    // Membuka koneksi ke database PostgreSQL
    db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
    if err != nil {
        logrus.WithError(err).Fatal("Failed to connect to database")
        helper.PanicIfError(err)
    }
    logrus.Info("Database connected successfully")

    // Menjalankan perintah SQL untuk membuat ekstensi uuid-ossp
    logrus.Info("Creating UUID extension if not exists...")
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
        logrus.WithError(err).Fatal("Failed to create uuid-ossp extension")
    }
    logrus.Info("UUID extension created or already exists")

    // Migrasi skema ke database
    logrus.Info("Running database migrations...")
    errMigrate := db.AutoMigrate(
        &schema.Category{},
        &schema.User{},
        &schema.Project{},
        &schema.ProjectItem{},
    )
    if errMigrate != nil {
        logrus.WithError(errMigrate).Fatal("Migration failed")
    }
    logrus.Info("Database migrations completed successfully")

    // RUN SEEDER
    logrus.Info("Running user seeder...")
    hashedPassword, errHash := helper.HashPassword(seedPassword)
    if errHash != nil {
        logrus.WithError(errHash).Fatal("Failed to hash password")
    }

    if errSeeder := db.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?) ON CONFLICT (email) DO NOTHING`, seedUsername, seedEmail, hashedPassword).Error; errSeeder != nil {
        logrus.WithError(errSeeder).Fatal("Failed to insert seed user")
    }
    logrus.WithFields(logrus.Fields{
        "username": seedUsername,
        "email":    seedEmail,
    }).Info("User seeding completed successfully")

    return db
}
