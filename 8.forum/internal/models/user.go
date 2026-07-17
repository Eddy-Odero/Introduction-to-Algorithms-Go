package models

import (
	"database/sql"
	"errors"
	"time"
)

// ErrDuplicateEmail is returned by CreateUser when the email is already
// registered, so the registration handler can turn it into a clean 409
// instead of a raw SQLite constraint error.
var ErrDuplicateEmail = errors.New("models: email already registered")

// ErrDuplicateUsername mirrors ErrDuplicateEmail for the username column.
var ErrDuplicateUsername = errors.New("models: username already taken")

// ErrNotFound is returned when a lookup finds no matching row.
var ErrNotFound = errors.New("models: not found")

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserStore struct {
	DB *sql.DB
}

// CreateUser inserts a new user row. The caller hashes the password
// before calling this — bcrypt hashing belongs in internal/auth (Phase 4),
// not in the data-access layer.
func (s *UserStore) CreateUser(email, username, passwordHash string) (int64, error) {
	res, err := s.DB.Exec(
		`INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)`,
		email, username, passwordHash,
	)
	if err != nil {
		msg := err.Error()
		switch {
		case msg == "UNIQUE constraint failed: users.email":
			return 0, ErrDuplicateEmail
		case msg == "UNIQUE constraint failed: users.username":
			return 0, ErrDuplicateUsername
		default:
			return 0, err
		}
	}
	return res.LastInsertId()
}

// GetUserByEmail looks up a user for login. Returns ErrNotFound if no row
// matches, so the login handler can show one generic "invalid
// credentials" message without revealing whether the email exists.
func (s *UserStore) GetUserByEmail(email string) (*User, error) {
	row := s.DB.QueryRow(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE email = ?`,
		email,
	)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetUserByID looks up a user by ID — used by the session middleware in
// Phase 4 to attach the current user to a request.
func (s *UserStore) GetUserByID(id int64) (*User, error) {
	row := s.DB.QueryRow(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE id = ?`,
		id,
	)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
