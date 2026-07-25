package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"forum/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a fresh in-memory database for each test
func setupTestDB(t *testing.T) *sql.DB {
	// Use :memory: for fast, isolated tests
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	// Enable foreign keys (critical for testing FK enforcement)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Read and execute schema
	schemaBytes, err := os.ReadFile(filepath.Join("..", "..", "internal", "db", "schema.sql"))
	if err != nil {
		// Try alternate path for test execution context
		schemaBytes, err = os.ReadFile("schema.sql")
		if err != nil {
			t.Fatalf("failed to read schema.sql: %v", err)
		}
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	// Seed categories (normally done in migrate())
	defaults := []string{"Tech", "Football", "Gaming", "Music", "Movies", "General"}
	for _, name := range defaults {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", name)
		if err != nil {
			t.Fatalf("failed to seed category %q: %v", name, err)
		}
	}

	return db
}

// Test 1: Database file is created + Tables exist
func TestSchemaCreation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Verify all expected tables exist
	tables := []string{"users", "sessions", "posts", "categories", "post_categories", "comments", "reactions"}
	for _, table := range tables {
		var name string
		err := db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
			table,
		).Scan(&name)
		if err == sql.ErrNoRows {
			t.Errorf("table %q does not exist", table)
		} else if err != nil {
			t.Errorf("failed to check table %q: %v", table, err)
		} else if name != table {
			t.Errorf("expected table %q, got %q", table, name)
		}
	}
	t.Log("✓ All tables exist")
}

func TestMigrateRenamesLegacyPasswordColumn(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		t.Fatalf("create legacy users table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO users (email, username, password) VALUES ('a@example.com', 'alice', 'hash')`); err != nil {
		t.Fatalf("seed legacy user: %v", err)
	}

	if err := migrate(db); err != nil {
		t.Fatalf("migrate legacy database: %v", err)
	}

	user, err := models.NewUserRepository(db).GetUserByEmail("a@example.com")
	if err != nil {
		t.Fatalf("read migrated user: %v", err)
	}
	if user.PasswordHash != "hash" {
		t.Errorf("expected preserved password hash, got %q", user.PasswordHash)
	}
}

// Test 2: Categories are seeded
func TestCategoriesSeeded(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	expected := []string{"Football", "Gaming", "General", "Movies", "Music", "Tech"}
	rows, err := db.Query("SELECT name FROM categories ORDER BY name")
	if err != nil {
		t.Fatalf("failed to query categories: %v", err)
	}
	defer rows.Close()

	var actual []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("failed to scan category: %v", err)
		}
		actual = append(actual, name)
	}

	if len(actual) != len(expected) {
		t.Errorf("expected %d categories, got %d", len(expected), len(actual))
	}
	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("expected category %q at position %d, got %q", exp, i, actual[i])
		}
	}
	t.Log("✓ Default categories seeded correctly")
}

// Test 3: User insertion works
func TestUserInsertion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := models.NewUserRepository(db)

	// Test successful creation
	id, err := repo.CreateUser("alice", "alice@example.com", "hashed_password_123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if id != 1 {
		t.Errorf("expected first user ID to be 1, got %d", id)
	}

	// Verify user exists
	user, err := repo.GetUserByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("failed to retrieve user by email: %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("expected username 'alice', got %q", user.Username)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got %q", user.Email)
	}
	if user.PasswordHash != "hashed_password_123" {
		t.Errorf("expected password hash to match")
	}

	// Verify lookup by username
	userByName, err := repo.GetUserByUsername("alice")
	if err != nil {
		t.Fatalf("failed to retrieve user by username: %v", err)
	}
	if userByName.ID != id {
		t.Errorf("expected user ID %d, got %d", id, userByName.ID)
	}

	// Verify lookup by ID
	userByID, err := repo.GetUserByID(id)
	if err != nil {
		t.Fatalf("failed to retrieve user by ID: %v", err)
	}
	if userByID.Username != "alice" {
		t.Errorf("expected username 'alice', got %q", userByID.Username)
	}

	t.Log("✓ User insertion and retrieval works")
}

// Test 4: Duplicate email is rejected
func TestDuplicateEmailRejected(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := models.NewUserRepository(db)

	// Create first user
	_, err := repo.CreateUser("alice", "alice@example.com", "hash1")
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	// Try duplicate email with different username
	_, err = repo.CreateUser("alice2", "alice@example.com", "hash2")
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}

	// Verify EmailExists returns true
	exists, err := repo.EmailExists("alice@example.com")
	if err != nil {
		t.Fatalf("EmailExists failed: %v", err)
	}
	if !exists {
		t.Error("expected EmailExists to return true for existing email")
	}

	// Verify EmailExists returns false for new email
	exists, err = repo.EmailExists("new@example.com")
	if err != nil {
		t.Fatalf("EmailExists failed: %v", err)
	}
	if exists {
		t.Error("expected EmailExists to return false for non-existent email")
	}

	t.Log("✓ Duplicate email properly rejected")
}

// Test 5: Duplicate username is rejected
func TestDuplicateUsernameRejected(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := models.NewUserRepository(db)

	// Create first user
	_, err := repo.CreateUser("alice", "alice@example.com", "hash1")
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	// Try duplicate username with different email
	_, err = repo.CreateUser("alice", "bob@example.com", "hash2")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}

	// Verify UsernameExists
	exists, err := repo.UsernameExists("alice")
	if err != nil {
		t.Fatalf("UsernameExists failed: %v", err)
	}
	if !exists {
		t.Error("expected UsernameExists to return true for existing username")
	}

	t.Log("✓ Duplicate username properly rejected")
}

// Test 6: Foreign keys are enforced
func TestForeignKeysEnforced(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := models.NewUserRepository(db)

	// Create a user
	userID, err := repo.CreateUser("alice", "alice@example.com", "hash")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Insert a post referencing the user (should succeed)
	_, err = db.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		userID, "Test Post", "Content",
	)
	if err != nil {
		t.Fatalf("failed to insert valid post: %v", err)
	}

	// Try inserting a post with non-existent user_id (should fail)
	_, err = db.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		9999, "Bad Post", "Content",
	)
	if err == nil {
		t.Fatal("expected foreign key error for invalid user_id, got nil")
	}

	// Try inserting a comment with non-existent post_id (should fail)
	_, err = db.Exec(
		"INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
		9999, userID, "Bad comment",
	)
	if err == nil {
		t.Fatal("expected foreign key error for invalid post_id, got nil")
	}

	// Verify CASCADE delete: deleting user should delete their posts
	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count posts: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 posts after user deletion (CASCADE), got %d", count)
	}

	t.Log("✓ Foreign keys enforced and CASCADE works")
}

// Test 7: Category repository methods
func TestCategoryRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := models.NewCategoryRepository(db)

	// Test GetAll
	cats, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(cats) != 6 {
		t.Errorf("expected 6 categories, got %d", len(cats))
	}

	// Test GetByName
	cat, err := repo.GetByName("Tech")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}
	if cat == nil {
		t.Fatal("expected to find 'Tech' category, got nil")
	}
	if cat.Name != "Tech" {
		t.Errorf("expected name 'Tech', got %q", cat.Name)
	}

	// Test GetByName for non-existent
	cat, err = repo.GetByName("NonExistent")
	if err != nil {
		t.Fatalf("GetByName failed for non-existent: %v", err)
	}
	if cat != nil {
		t.Error("expected nil for non-existent category")
	}

	// Test LinkToPost - create a post first
	userRepo := models.NewUserRepository(db)
	userID, _ := userRepo.CreateUser("bob", "bob@example.com", "hash")
	postRes, _ := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, "Post", "Content")
	postID, _ := postRes.LastInsertId()

	// Link categories
	err = repo.LinkToPost(int(postID), []string{"Tech", "Gaming"})
	if err != nil {
		t.Fatalf("LinkToPost failed: %v", err)
	}

	// Verify via GetByPostID
	names, err := repo.GetByPostID(int(postID))
	if err != nil {
		t.Fatalf("GetByPostID failed: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 categories linked, got %d", len(names))
	}
	if names[0] != "Gaming" || names[1] != "Tech" {
		t.Errorf("expected [Gaming, Tech], got %v", names)
	}

	// Test auto-create unknown category
	err = repo.LinkToPost(int(postID), []string{"NewCategory"})
	if err != nil {
		t.Fatalf("LinkToPost with new category failed: %v", err)
	}

	// Verify new category was created
	cat, err = repo.GetByName("NewCategory")
	if err != nil {
		t.Fatalf("GetByName for auto-created category failed: %v", err)
	}
	if cat == nil {
		t.Fatal("expected auto-created category to exist")
	}

	t.Log("✓ Category repository works correctly")
}

// Test 8: Unique constraints on reactions
func TestReactionUniqueConstraints(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := models.NewUserRepository(db)
	uid1, _ := userRepo.CreateUser("u1", "u1@example.com", "h1")
	uid2, _ := userRepo.CreateUser("u2", "u2@example.com", "h2")

	postRes, _ := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", uid1, "P1", "C1")
	pid, _ := postRes.LastInsertId()

	// User 1 likes post 1
	_, err := db.Exec("INSERT INTO reactions (user_id, post_id, type) VALUES (?, ?, 'like')", uid1, pid)
	if err != nil {
		t.Fatalf("first reaction insert failed: %v", err)
	}

	// User 1 tries to like same post again (should fail - UNIQUE constraint)
	_, err = db.Exec("INSERT INTO reactions (user_id, post_id, type) VALUES (?, ?, 'like')", uid1, pid)
	if err == nil {
		t.Fatal("expected unique constraint violation for duplicate reaction")
	}

	// User 2 likes same post (should succeed - different user)
	_, err = db.Exec("INSERT INTO reactions (user_id, post_id, type) VALUES (?, ?, 'like')", uid2, pid)
	if err != nil {
		t.Fatalf("different user reaction failed: %v", err)
	}

	t.Log("✓ Reaction unique constraints work")
}
