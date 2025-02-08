package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"appgents/internal/core/auth"
	"appgents/internal/core/email"
	"appgents/internal/platform/config"
	"appgents/internal/platform/database"
	authdb "appgents/internal/platform/database/auth"
	"appgents/internal/platform/database/migrations"
	httpserver "appgents/internal/platform/http"
	authhandler "appgents/internal/platform/http/handler/auth"
	"appgents/internal/platform/http/middleware"
	"appgents/internal/platform/logging"
)

func main() {
	// Initialize logger
	logger := logging.New()
	logger.Info("starting_server", nil)

	// Load config
	config, err := config.Load("config.json")
	if err != nil {
		logger.Error("Failed to load config", err, nil)
		os.Exit(1)
	}

	// Initialize email service with logger
	emailService, err := email.NewEmailService(config.Email, logger)
	if err != nil {
		logger.Error("Failed to initialize email service", err, nil)
		os.Exit(1)
	}

	sender := email.NewSender(emailService, config.Server.BaseURL)

	// Initialize components
	// eventBus := bus.New()
	// appgentRegistry := registry.New()

	// Set up database
	db, err := database.New(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set up repositories and services
	userRepo := authdb.NewRepository(db)
	authService := auth.NewService(userRepo, sender)

	// Create handlers
	authHandler := authhandler.New(logger, authService)

	// Create router
	mux := http.NewServeMux()
	setupRoutes(mux, logger, authService)

	// Add logging middleware
	mux.Use(middleware.RequestLogger(logger))

	// Configure and create server
	serverConfig := httpserver.Config{
		Port:         "8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	server := httpserver.NewServer(mux, serverConfig)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	// Start server
	go func() {
		log.Printf("Server listening on :%s", serverConfig.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
