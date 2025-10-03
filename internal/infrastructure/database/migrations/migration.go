// internal/infrastructure/database/migrations/migrate.go
package migrations

import (
	"fmt"
	"log"
	"sort"

	"gorm.io/gorm"
)

// Migration struct untuk mendefinisikan migration
type Migration struct {
	Version int
	Name    string
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error
}

// Slice untuk menyimpan semua migrations
var migrations []Migration

// Register function untuk menambahkan migration ke slice
// ⬇️ INI DIA FUNCTION REGISTER
func Register(m Migration) {
	migrations = append(migrations, m)
}

// MigrationHistory untuk track migration yang sudah dijalankan
type MigrationHistory struct {
	ID        uint   `gorm:"primaryKey"`
	Version   int    `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	AppliedAt int64  `gorm:"autoCreateTime"`
}

// RunMigrations menjalankan semua pending migrations
func RunMigrations(db *gorm.DB) error {
	// Create migration history table jika belum ada
	if err := db.AutoMigrate(&MigrationHistory{}); err != nil {
		return fmt.Errorf("failed to create migration history table: %w", err)
	}

	// Sort migrations berdasarkan version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Get migrations yang sudah dijalankan
	var applied []MigrationHistory
	db.Order("version").Find(&applied)
	appliedMap := make(map[int]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	// Run pending migrations
	for _, migration := range migrations {
		// Skip jika sudah dijalankan
		if appliedMap[migration.Version] {
			log.Printf("⏭️  Skipping migration %d: %s (already applied)", migration.Version, migration.Name)
			continue
		}

		log.Printf("⬆️  Running migration %d: %s", migration.Version, migration.Name)

		// Jalankan migration dalam transaction
		err := db.Transaction(func(tx *gorm.DB) error {
			// Execute UP function
			if err := migration.Up(tx); err != nil {
				return err
			}

			// Catat ke migration history
			return tx.Create(&MigrationHistory{
				Version: migration.Version,
				Name:    migration.Name,
			}).Error
		})

		if err != nil {
			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		log.Printf("✅ Migration %d completed: %s", migration.Version, migration.Name)
	}

	log.Println("🎉 All migrations completed successfully!")
	return nil
}

// Rollback untuk rollback N migrations terakhir
func Rollback(db *gorm.DB, steps int) error {
	// Get applied migrations (newest first)
	var applied []MigrationHistory
	db.Order("version DESC").Limit(steps).Find(&applied)

	if len(applied) == 0 {
		log.Println("ℹ️  No migrations to rollback")
		return nil
	}

	// Rollback migrations
	for _, history := range applied {
		// Find migration definition
		var migration *Migration
		for _, m := range migrations {
			if m.Version == history.Version {
				migration = &m
				break
			}
		}

		if migration == nil {
			return fmt.Errorf("migration %d not found in code", history.Version)
		}

		log.Printf("⬇️  Rolling back migration %d: %s", migration.Version, migration.Name)

		// Rollback dalam transaction
		err := db.Transaction(func(tx *gorm.DB) error {
			// Execute DOWN function
			if err := migration.Down(tx); err != nil {
				return err
			}

			// Hapus dari migration history
			return tx.Where("version = ?", migration.Version).Delete(&MigrationHistory{}).Error
		})

		if err != nil {
			return fmt.Errorf("rollback migration %d failed: %w", migration.Version, err)
		}

		log.Printf("✅ Rolled back migration %d: %s", migration.Version, migration.Name)
	}

	return nil
}

// Status menampilkan status semua migrations
func Status(db *gorm.DB) error {
	// Create table if not exists
	db.AutoMigrate(&MigrationHistory{})

	// Get applied migrations
	var applied []MigrationHistory
	db.Order("version").Find(&applied)
	appliedMap := make(map[int]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	// Sort migrations
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	log.Println("\n📋 Migration Status:")
	log.Println("-------------------")
	for _, m := range migrations {
		status := "❌ Pending"
		if appliedMap[m.Version] {
			status = "✅ Applied"
		}
		log.Printf("Version %d: %s - %s", m.Version, m.Name, status)
	}
	log.Println()

	return nil
}
