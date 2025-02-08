package auth

import (
	"context"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
	"github.com/gee-m/sidekick/web/templates/auth"
	"github.com/go-chi/chi/v5"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", h.handleLoginPage)
	r.Post("/login", h.handleLogin)
	r.Post("/check-email", h.CheckEmail)

	return r
}

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	component := auth.LoginPage(auth.LoginProps{})
	component.Render(r.Context(), w)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// To be implemented
	// This will handle the actual login logic later
	component := auth.LoginPage(auth.LoginProps{
		ErrorMessage: "Login functionality coming soon!",
	})
	component.Render(r.Context(), w)
}

func (h *Handler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// First validate email format
	if !emailRegex.MatchString(email) {
		h.render(r.Context(), w, auth.ActionButtons(false))
		return
	}

	// Check if email exists
	exists, err := h.service.EmailExists(r.Context(), email)
	if err != nil {
		// On error, just render without signup button
		h.render(r.Context(), w, auth.ActionButtons(false))
		return
	}

	// Show signup button only if email doesn't exist
	h.render(r.Context(), w, auth.ActionButtons(!exists))
}

func (h *Handler) render(ctx context.Context, w http.ResponseWriter, comp templ.Component) {
	comp.Render(ctx, w)
}
