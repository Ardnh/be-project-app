// internal/infrastructure/database/migrations/006_create_project_todolists_table.go
package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 6,
		Name:    "create_project_todolists_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS project_todolists (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					project_id UUID NOT NULL,
					name VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
				);

				CREATE INDEX idx_project_todolists_project_id ON project_todolists(project_id);
				CREATE INDEX idx_project_todolists_deleted_at ON project_todolists(deleted_at);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS project_todolists CASCADE;").Error
		},
	})
}
