package database

import (
	"database/sql"
	"log"
)

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_roles_table",
			sql: `
				CREATE TABLE IF NOT EXISTS roles (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
				CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
			`,
		},
		{
			name: "create_users_table",
			sql: `
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					email VARCHAR(255) NOT NULL UNIQUE,
					password VARCHAR(255) NOT NULL,
					role_id INTEGER NOT NULL DEFAULT 1 REFERENCES roles(id),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
			`,
		},
		{
			name: "insert_default_roles",
			sql: `
				INSERT INTO roles (name) VALUES ('user') ON CONFLICT (name) DO NOTHING;
				INSERT INTO roles (name) VALUES ('admin') ON CONFLICT (name) DO NOTHING;
			`,
		},
		{
			name: "create_permissions_table",
			sql: `
				CREATE TABLE IF NOT EXISTS permissions (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
				CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);
			`,
		},
		{
			name: "create_role_permissions_table",
			sql: `
				CREATE TABLE IF NOT EXISTS role_permissions (
					role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
					permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
					PRIMARY KEY (role_id, permission_id)
				);
				CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
				CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
			`,
		},
	}

	for _, migration := range migrations {
		log.Printf("Running migration: %s\n", migration.name)
		if _, err := db.Exec(migration.sql); err != nil {
			return err
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}
