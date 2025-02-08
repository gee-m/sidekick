package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gee-m/sidekick/internal/core/auth"
	"github.com/gee-m/sidekick/internal/platform/database"
)

const (
	defaultPort     = "3000"
	shutdownTimeout = 10 * time.Second
)

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Get database configuration
	dbConfig := database.NewConfig()

	// Run migrations
	if err := database.RunMigrations(dbConfig.URL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database connection
	db, err := database.Connect(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services and handlers
	authService := auth.NewService(db)
	sessions := auth.NewSessionManager(db)

	// Setup routes with auth middleware
	router := setupRoutes(authService, sessions)

	// Create server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
