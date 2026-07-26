package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrent read/write performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Set busy timeout (5 seconds) to prevent "database is locked" errors
	if _, err := db.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	// Run schema migration
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	DB = db
	return db, nil
}

// migrate reads schema.sql and executes it
func migrate(db *sql.DB) error {
	if err := migratePasswordColumn(db); err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(0)
	schemaPath := filepath.Join(filepath.Dir(filename), "schema.sql")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Seed default categories if table is empty
	return seedCategories(db)
}

// migratePasswordColumn preserves accounts created before the password_hash
// column name was standardized in Phase 4.
func migratePasswordColumn(db *sql.DB) error {
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		return fmt.Errorf("inspect users table: %w", err)
	}
	defer rows.Close()

	var hasPassword, hasPasswordHash bool
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull, primaryKey int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return fmt.Errorf("read users column: %w", err)
		}
		switch name {
		case "password":
			hasPassword = true
		case "password_hash":
			hasPasswordHash = true
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate users columns: %w", err)
	}
	if hasPassword && !hasPasswordHash {
		if _, err := db.Exec("ALTER TABLE users RENAME COLUMN password TO password_hash"); err != nil {
			return fmt.Errorf("rename users.password to password_hash: %w", err)
		}
	}
	return nil
}

// seedCategories inserts default forum categories
func seedCategories(db *sql.DB) error {
	defaults := []string{"Tech", "Football", "Gaming", "Music", "Movies", "General"}

	for _, name := range defaults {
		_, err := db.Exec(
			"INSERT OR IGNORE INTO categories (name) VALUES (?)",
			name,
		)
		if err != nil {
			return fmt.Errorf("failed to seed category %q: %w", name, err)
		}
	}
	return nil
}

// Close cleanly shuts down the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
