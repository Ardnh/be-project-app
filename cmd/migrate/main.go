package main

import (
	"flag"
	"log"
	"os"

	"github.com/Ardnh/be-project-app/internal/config"
	"github.com/Ardnh/be-project-app/internal/infrastructure/database/migrations"
	"github.com/Ardnh/be-project-app/internal/infrastructure/database/postgresql"

	_ "github.com/Ardnh/be-project-app/internal/infrastructure/database/migrations"
)

func main() {

	// Load config
	cfg := config.LoadConfig()

	// Parse command line flags
	action := flag.String("action", "up", "Migration action: up, down, status")
	steps := flag.Int("steps", 1, "Number of migrations to rollback (for down action)")
	flag.Parse()

	// Connect to database
	db, err := postgresql.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Run migration action
	switch *action {
	case "up":
		if err := migrations.RunMigrations(db); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}

	case "down":
		if err := migrations.Rollback(db, *steps); err != nil {
			log.Fatalf("❌ Rollback failed: %v", err)
		}

	case "status":
		if err := migrations.Status(db); err != nil {
			log.Fatalf("❌ Status check failed: %v", err)
		}

	default:
		log.Fatalf("❌ Unknown action: %s. Use: up, down, or status", *action)
	}

	os.Exit(0)
}
