// Package db owns the single *sql.DB handle for the forum and knows how
// to initialize the schema. Every other package gets its connection by
// calling db.Open, not by importing sqlite3 directly.
package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

// Open creates (if needed) the SQLite file at path, applies pragmas that
// matter for a web server workload, runs the schema, and returns the
// handle.
func Open(path string) (*sql.DB, error) {
	// _foreign_keys=on   -> enforce FK constraints (off by default in SQLite)
	// _journal_mode=WAL  -> readers don't block writers under concurrent
	//                       requests
	// _busy_timeout=5000 -> wait up to 5s on a locked DB instead of failing
	//                       immediately with "database is locked"
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000", path)

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: open %s: %w", path, err)
	}

	// SQLite only supports one writer at a time; capping the pool at a
	// single connection avoids "database is locked" errors from Go's
	// connection pool trying to write concurrently, rather than chasing
	// that error down later.
	conn.SetMaxOpenConns(1)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	if err := migrate(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("db: migrate: %w", err)
	}

	return conn, nil
}

// migrate runs schema.sql. CREATE TABLE/INDEX use IF NOT EXISTS and seed
// data uses INSERT OR IGNORE, so this is safe to run on every startup.
func migrate(conn *sql.DB) error {
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema.sql: %w", err)
	}
	if _, err := conn.Exec(string(schema)); err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}
	return nil
}
