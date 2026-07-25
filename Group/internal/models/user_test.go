package models_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"forum/internal/db"
	"forum/internal/models"
)

// openTestDB creates a fresh SQLite file in a temp dir per test, so tests
// never share state or touch the real forum.db.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	conn, err := db.Open(path)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestCreateAndGetUser(t *testing.T) {
	conn := openTestDB(t)
	users := &models.UserStore{DB: conn}

	id, err := users.CreateUser("eddy@example.com", "eddy_k", "hashed-password")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero user id")
	}

	got, err := users.GetUserByEmail("eddy@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if got.Username != "eddy_k" {
		t.Errorf("username = %q, want %q", got.Username, "eddy_k")
	}
}

func TestDuplicateEmailRejected(t *testing.T) {
	conn := openTestDB(t)
	users := &models.UserStore{DB: conn}

	if _, err := users.CreateUser("dup@example.com", "user_one", "hash1"); err != nil {
		t.Fatalf("first CreateUser: %v", err)
	}

	_, err := users.CreateUser("dup@example.com", "user_two", "hash2")
	if err != models.ErrDuplicateEmail {
		t.Errorf("got err = %v, want ErrDuplicateEmail", err)
	}
}

func TestDuplicateUsernameRejected(t *testing.T) {
	conn := openTestDB(t)
	users := &models.UserStore{DB: conn}

	if _, err := users.CreateUser("a@example.com", "sameuser", "hash1"); err != nil {
		t.Fatalf("first CreateUser: %v", err)
	}

	_, err := users.CreateUser("b@example.com", "sameuser", "hash2")
	if err != models.ErrDuplicateUsername {
		t.Errorf("got err = %v, want ErrDuplicateUsername", err)
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	conn := openTestDB(t)
	users := &models.UserStore{DB: conn}

	_, err := users.GetUserByEmail("nobody@example.com")
	if err != models.ErrNotFound {
		t.Errorf("got err = %v, want ErrNotFound", err)
	}
}
