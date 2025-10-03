// internal/infrastructure/database/migrations/002_create_categories_table.go
package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 2,
		Name:    "create_categories_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS categories (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(100) UNIQUE NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				);

				CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS categories CASCADE;").Error
		},
	})
}
