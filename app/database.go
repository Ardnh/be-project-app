package app

import (
	"os"
	"project-app/helper"
	"project-app/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {

	postgresDsn := os.Getenv("APP_DSN")
	db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})

	helper.PanicIfError(err)

	db.AutoMigrate(
		&schema.Categories{},
		&schema.Users{},
		&schema.Projects{},
		&schema.ProjectItem{},
	)

	return db
}
