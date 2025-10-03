package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: 1,
		Name:    "create_users_table",
		Up: func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					username VARCHAR(100) UNIQUE NOT NULL,
					email VARCHAR(255) UNIQUE NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				);

				CREATE INDEX idx_users_deleted_at ON users(deleted_at);
				CREATE INDEX idx_users_username ON users(username);
				CREATE INDEX idx_users_email ON users(email);
			`).Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS users CASCADE;").Error
		},
	})
}
