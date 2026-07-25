package models

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupCommentTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable fk: %v", err)
	}
	schema, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("exec schema: %v", err)
	}
	return db
}

func seedCommentUser(t *testing.T, db *sql.DB) int64 {
	res, err := db.Exec(
		"INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)",
		"t@test.com", "tester", "hash",
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}

func seedCommentPost(t *testing.T, db *sql.DB, userID int64) int64 {
	res, err := db.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		userID, "Title", "Content",
	)
	if err != nil {
		t.Fatalf("seed post: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}

func TestCreateComment(t *testing.T) {
	db := setupCommentTestDB(t)
	defer db.Close()
	repo := NewCommentRepository(db)
	user := seedCommentUser(t, db)
	post := seedCommentPost(t, db, user)

	c, err := repo.CreateComment(post, user, "Nice post!")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if c.Content != "Nice post!" {
		t.Errorf("content: got %q", c.Content)
	}
	if c.Username != "tester" {
		t.Errorf("username: got %q", c.Username)
	}
}

func TestCreateCommentEmpty(t *testing.T) {
	db := setupCommentTestDB(t)
	defer db.Close()
	repo := NewCommentRepository(db)
	user := seedCommentUser(t, db)
	post := seedCommentPost(t, db, user)

	_, err := repo.CreateComment(post, user, "")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestGetCommentByID(t *testing.T) {
	db := setupCommentTestDB(t)
	defer db.Close()
	repo := NewCommentRepository(db)
	user := seedCommentUser(t, db)
	post := seedCommentPost(t, db, user)

	created, _ := repo.CreateComment(post, user, "Find me")
	found, err := repo.GetCommentByID(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if found.Content != "Find me" {
		t.Errorf("content: got %q", found.Content)
	}
}

func TestGetCommentsByPost(t *testing.T) {
	db := setupCommentTestDB(t)
	defer db.Close()
	repo := NewCommentRepository(db)
	user := seedCommentUser(t, db)
	post := seedCommentPost(t, db, user)

	repo.CreateComment(post, user, "First")
	repo.CreateComment(post, user, "Second")

	comments, err := repo.GetCommentsByPost(post)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("count: got %d", len(comments))
	}
	if comments[0].Content != "First" {
		t.Errorf("order: got %q first", comments[0].Content)
	}
}

func TestCountCommentsByPost(t *testing.T) {
	db := setupCommentTestDB(t)
	defer db.Close()
	repo := NewCommentRepository(db)
	user := seedCommentUser(t, db)
	post := seedCommentPost(t, db, user)

	count, _ := repo.CountCommentsByPost(post)
	if count != 0 {
		t.Errorf("initial: got %d", count)
	}

	repo.CreateComment(post, user, "One")
	repo.CreateComment(post, user, "Two")

	count, _ = repo.CountCommentsByPost(post)
	if count != 2 {
		t.Errorf("after: got %d", count)
	}
}

func TestDeleteComment(t *testing.T) {
	db := setupCommentTestDB(t)
	defer db.Close()
	repo := NewCommentRepository(db)
	user := seedCommentUser(t, db)
	post := seedCommentPost(t, db, user)

	c, _ := repo.CreateComment(post, user, "Delete me")
	if err := repo.DeleteComment(c.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := repo.GetCommentByID(c.ID)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
}
