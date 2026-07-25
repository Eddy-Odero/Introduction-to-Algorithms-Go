package models

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database with the users table
func setupTestDB(t *testing.T) *sql.DB {
	// :memory: creates a temporary database that exists only during the test
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the users table (same schema as your teammate's schema.sql)
	createTableSQL := `CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

// cleanup closes the database after each test
func cleanup(db *sql.DB) {
	db.Close()
}

// TestCreateUser verifies that CreateUser inserts a new user correctly
func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Insert a user
	id, err := repo.CreateUser("testuser", "test@forum.com", "hashedpassword123")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Verify we got an ID back
	if id == 0 {
		t.Error("Expected non-zero ID, got 0")
	}
}

// TestGetUserByEmail verifies retrieving a user by email
func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// First, create a user
	_, err := repo.CreateUser("alice", "alice@forum.com", "hash456")
	if err != nil {
		t.Fatalf("Setup: CreateUser failed: %v", err)
	}

	// Now retrieve by email
	user, err := repo.GetUserByEmail("alice@forum.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	// Check all fields
	if user.Username != "alice" {
		t.Errorf("Expected username 'alice', got '%s'", user.Username)
	}
	if user.Email != "alice@forum.com" {
		t.Errorf("Expected email 'alice@forum.com', got '%s'", user.Email)
	}
	if user.PasswordHash != "hash456" {
		t.Errorf("Expected password hash 'hash456', got '%s'", user.PasswordHash)
	}
}

// TestGetUserByEmail_NotFound verifies error when email doesn't exist
func TestGetUserByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Try to get a user that was never created
	_, err := repo.GetUserByEmail("nobody@forum.com")
	if err == nil {
		t.Error("Expected an error for non-existent email, got nil")
	}
}

// TestGetUserByUsername verifies retrieving a user by username
func TestGetUserByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	repo.CreateUser("bob", "bob@forum.com", "hash789")

	user, err := repo.GetUserByUsername("bob")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if user.Email != "bob@forum.com" {
		t.Errorf("Expected email 'bob@forum.com', got '%s'", user.Email)
	}
}

// TestGetUserByID verifies retrieving a user by ID
func TestGetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Create user and capture the ID
	id, _ := repo.CreateUser("charlie", "charlie@forum.com", "hash000")

	user, err := repo.GetUserByID(id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if user.Username != "charlie" {
		t.Errorf("Expected username 'charlie', got '%s'", user.Username)
	}
}

// TestEmailExists verifies email existence check
func TestEmailExists(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Check before creating user
	exists, _ := repo.EmailExists("dave@forum.com")
	if exists {
		t.Error("Expected email to NOT exist before creation")
	}

	// Create user
	repo.CreateUser("dave", "dave@forum.com", "hash111")

	// Check after creating user
	exists, _ = repo.EmailExists("dave@forum.com")
	if !exists {
		t.Error("Expected email to exist after creation")
	}
}

// TestUsernameExists verifies username existence check
func TestUsernameExists(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Check before creating
	exists, _ := repo.UsernameExists("eve")
	if exists {
		t.Error("Expected username to NOT exist before creation")
	}

	// Create user
	repo.CreateUser("eve", "eve@forum.com", "hash222")

	// Check after creating
	exists, _ = repo.UsernameExists("eve")
	if !exists {
		t.Error("Expected username to exist after creation")
	}
}

// TestDuplicateEmail verifies that duplicate emails are rejected
func TestDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Create first user
	_, err := repo.CreateUser("user1", "same@forum.com", "hash1")
	if err != nil {
		t.Fatalf("First CreateUser failed: %v", err)
	}

	// Try to create second user with same email
	_, err = repo.CreateUser("user2", "same@forum.com", "hash2")
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}

// TestDuplicateUsername verifies that duplicate usernames are rejected
func TestDuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer cleanup(db)

	repo := NewUserRepository(db)

	// Create first user
	repo.CreateUser("uniqueuser", "one@forum.com", "hash1")

	// Try to create second user with same username
	_, err := repo.CreateUser("uniqueuser", "two@forum.com", "hash2")
	if err == nil {
		t.Error("Expected error for duplicate username, got nil")
	}
}
