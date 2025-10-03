// internal/infrastructure/database/migrations/007_create_project_todolist_items_table.go
package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 7,
		Name:    "create_project_todolist_items_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS project_todolist_items (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					project_todolist_id UUID NOT NULL,
					name VARCHAR(255) NOT NULL,
					status VARCHAR(50) DEFAULT 'pending',
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_todolist_id) REFERENCES project_todolists(id) ON DELETE CASCADE
				);

				CREATE INDEX idx_project_todolist_items_project_todolist_id ON project_todolist_items(project_todolist_id);
				CREATE INDEX idx_project_todolist_items_deleted_at ON project_todolist_items(deleted_at);
				CREATE INDEX idx_project_todolist_items_status ON project_todolist_items(status);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS project_todolist_items CASCADE;").Error
		},
	})
}
