package main

import (
	"forum/internal/db"
	"forum/internal/middleware"
	"forum/internal/router"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := envOrDefault("FORUM_ADDR", ":8080")
	uiDir := envOrDefault("FORUM_UI_DIR", "web")
	dbPath := envOrDefault("FORUM_DB_PATH", "forum.db")

	_, err := db.Init(dbPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Database initialized", "path", dbPath)

	mux := router.SetupRouter(uiDir)
	
	// Wrap with auth middleware (Phase 4)
	wrappedRouter := middleware.AuthMiddleware(mux)
	wrappedRouter = middleware.LoggingMiddleware(wrappedRouter)

	slog.Info("Server successfully started", "port", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           wrappedRouter,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Critical server shutdown error", "error", err)
		os.Exit(1)
	}
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}