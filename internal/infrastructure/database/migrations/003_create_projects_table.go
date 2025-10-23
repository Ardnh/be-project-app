// internal/infrastructure/database/migrations/003_create_projects_table.go
package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 3,
		Name:    "create_projects_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS projects (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID NOT NULL,
					name VARCHAR(255) NOT NULL,
					budget DECIMAL(15,2) DEFAULT 0,
					category_name VARCHAR(255) NOT NULL,
					is_completed BOOLEAN DEFAULT false,
					start_date DATE,
					end_date DATE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);
				CREATE INDEX idx_projects_user_id ON projects(user_id);
				CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS projects CASCADE;").Error
		},
	})
}
