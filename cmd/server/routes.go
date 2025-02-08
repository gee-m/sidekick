package main

import (
	"net/http"

	authn "appgents/internal/core/auth"
	"appgents/internal/platform/http/handler/auth"
	"appgents/internal/platform/logging"
)

func setupRoutes(mux *http.ServeMux, logger *logging.Logger, authService *authn.Service) {
	// Auth routes
	authHandler := auth.New(logger, authService)
	mux.Handle("/", authHandler)
	mux.Handle("/auth/login", authHandler)

	// Health check
	mux.HandleFunc("/health", handleHealth)

	// Add more routes here as we build them
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
