package models

import (
	"database/sql"
	"errors"
	"time"
)

// User represents a registered forum user
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// UserRepository provides database operations for users
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository with the given database connection
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser inserts a new user into the database
// Returns the new user's ID or an error if email/username already exists
func (ur *UserRepository) CreateUser(username, email, passwordHash string) (int64, error) {
	result, err := ur.DB.Exec(
		`INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`,
		username, email, passwordHash,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByEmail retrieves a user by their email address
// Returns sql.ErrNoRows if not found
func (ur *UserRepository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := ur.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username
// Returns sql.ErrNoRows if not found
func (ur *UserRepository) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := ur.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID
// Returns sql.ErrNoRows if not found
func (ur *UserRepository) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := ur.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// EmailExists checks if an email is already registered
func (ur *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	err := ur.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`,
		email,
	).Scan(&exists)
	return exists, err
}

// UsernameExists checks if a username is already taken
func (ur *UserRepository) UsernameExists(username string) (bool, error) {
	var exists bool
	err := ur.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`,
		username,
	).Scan(&exists)
	return exists, err
}

// Common errors
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
)
