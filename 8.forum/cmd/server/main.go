package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum/internal/db"
	"forum/internal/middleware"
	"forum/internal/router"
	"forum/internal/view"
)

const addr = ":8080"

func main() {
	if err := view.Load("web/templates"); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	dbPath := os.Getenv("FORUM_DB_PATH")
	if dbPath == "" {
		dbPath = "./forum.db"
	}
	conn, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer conn.Close()
	log.Printf("database ready at %s", dbPath)

	handler := middleware.Logging(router.New())

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in a goroutine so main can block on the shutdown signal.
	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until we receive Ctrl+C or a termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped cleanly")
}