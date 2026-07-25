package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	SessionCookieName = "forum_session"
	SessionDuration   = 24 * time.Hour
)

// Session represents a user session stored in the database
type Session struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// SessionRepository provides database operations for sessions
type SessionRepository struct {
	DB *sql.DB
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (sr *SessionRepository) CreateSession(userID int) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(SessionDuration)

	// Enforce single active session: invalidate all previous sessions for this user
	_, err := sr.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return "", fmt.Errorf("delete old sessions failed: %w", err)
	}

	_, err = sr.DB.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		userID, token, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("create session failed: %w", err)
	}
	return token, nil
}
// GetSessionByToken retrieves a valid (non-expired) session by token
func (sr *SessionRepository) GetSessionByToken(token string) (*Session, error) {
	session := &Session{}
	err := sr.DB.QueryRow(
		"SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = ? AND expires_at > ?",
		token, time.Now(),
	).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get session failed: %w", err)
	}
	return session, nil
}

// DeleteSession removes a session by token (logout)
func (sr *SessionRepository) DeleteSession(token string) error {
	_, err := sr.DB.Exec("DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		return fmt.Errorf("delete session failed: %w", err)
	}
	return nil
}

// DeleteExpiredSessions cleans up expired sessions
func (sr *SessionRepository) DeleteExpiredSessions() error {
	_, err := sr.DB.Exec("DELETE FROM sessions WHERE expires_at <= ?", time.Now())
	if err != nil {
		return fmt.Errorf("delete expired sessions failed: %w", err)
	}
	return nil
}

// SetSessionCookie sets the session cookie on the response
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(SessionDuration.Seconds()),
	})
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

// GetSessionTokenFromCookie reads the session token from the request cookie
func GetSessionTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
