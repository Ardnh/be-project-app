// internal/infrastructure/database/migrations/004_create_project_expenses_table.go
package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 4,
		Name:    "create_project_expenses_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS project_expenses (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					project_id UUID NOT NULL,
					name VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
				);

				CREATE INDEX idx_project_expenses_project_id ON project_expenses(project_id);
				CREATE INDEX idx_project_expenses_deleted_at ON project_expenses(deleted_at);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS project_expenses CASCADE;").Error
		},
	})
}
