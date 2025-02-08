package main

import (
	"github.com/gee-m/sidekick/internal/core/auth"
	"github.com/gee-m/sidekick/internal/core/dashboard"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func setupRoutes(authService *auth.Service, sessions *auth.SessionManager) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Root redirect to dashboard
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	// Auth routes
	authHandler := auth.NewHandler(authService, sessions)
	r.Mount("/auth", authHandler.Routes())

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(sessions.AuthMiddleware)

		// Dashboard routes
		dashboardHandler := dashboard.NewHandler(authService)
		r.Mount("/dashboard", dashboardHandler.Routes())
	})

	return r
}
