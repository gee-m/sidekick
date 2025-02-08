package main

import (
	"github.com/gee-m/sidekick/internal/core/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRoutes(authService *auth.Service) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth routes
	authHandler := auth.NewHandler(authService)
	r.Mount("/auth", authHandler.Routes())

	return r
}
