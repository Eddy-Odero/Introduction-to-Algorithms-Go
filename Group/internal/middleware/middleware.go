package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"forum/internal/auth"
	"forum/internal/db"
	"forum/internal/models"
)

// UserContextKey is the key for storing user in request context
type UserContextKey struct{}

// UserContext holds user data in context
type UserContext struct {
	ID       int64
	Username string
	Email    string
	LoggedIn bool
}

// AuthMiddleware checks session cookie and loads user into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := auth.GetSessionTokenFromCookie(r)
		if token == "" {
			// No cookie, continue as guest
			ctx := context.WithValue(r.Context(), UserContextKey{}, UserContext{LoggedIn: false})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Look up session
		sessionRepo := auth.NewSessionRepository(db.DB)
		session, err := sessionRepo.GetSessionByToken(token)
		if err != nil || session == nil {
			// Invalid/expired session, clear cookie
			auth.ClearSessionCookie(w)
			ctx := context.WithValue(r.Context(), UserContextKey{}, UserContext{LoggedIn: false})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Load user
		userRepo := models.NewUserRepository(db.DB)
		user, err := userRepo.GetUserByID(int64(session.UserID))
		if err != nil || user == nil {
			ctx := context.WithValue(r.Context(), UserContextKey{}, UserContext{LoggedIn: false})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Attach user to context
		ctx := context.WithValue(r.Context(), UserContextKey{}, UserContext{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			LoggedIn: true,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves user from request context
func GetUserFromContext(r *http.Request) UserContext {
	user, ok := r.Context().Value(UserContextKey{}).(UserContext)
	if !ok {
		return UserContext{LoggedIn: false}
	}
	return user
}

// RequireAuth wraps a handler to require authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if !user.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// CustomNotFoundHandler handles 404 errors
func CustomNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, "web/templates/404.html")
}

// LoggingMiddleware logs each HTTP request with method, path, status, and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", duration,
			"remote_addr", r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}