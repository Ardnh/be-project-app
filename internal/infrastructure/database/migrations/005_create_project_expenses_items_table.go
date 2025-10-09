// internal/infrastructure/database/migrations/005_create_project_expense_items_table.go
package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 5,
		Name:    "create_project_expense_items_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS project_expense_items (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					project_expense_id UUID NOT NULL,
					amount DECIMAL(15,2) NOT NULL,
					name VARCHAR(255) NOT NULL,
					category_name VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_expense_id) REFERENCES project_expenses(id) ON DELETE CASCADE
				);

				CREATE INDEX idx_project_expense_items_project_expense_id ON project_expense_items(project_expense_id);
				CREATE INDEX idx_project_expense_items_deleted_at ON project_expense_items(deleted_at);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS project_expense_items CASCADE;").Error
		},
	})
}
