package auth

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	"github.com/gee-m/sidekick/web/templates/auth"
	"github.com/go-chi/chi/v5"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

const MinPasswordLength = 6

type Handler struct {
	service  *Service
	sessions *SessionManager
}

func NewHandler(service *Service, sessions *SessionManager) *Handler {
	return &Handler{
		service:  service,
		sessions: sessions,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", h.handleLoginPage)
	r.Post("/login", h.handleLogin)
	r.Post("/check-email", h.CheckEmail)
	r.Post("/signup", h.handleSignup)
	r.Post("/logout", h.handleLogout)             // Add logout route
	r.Get("/auth/verify/{token}", h.handleVerify) // Add verify route

	return r
}

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Try to get session cookie first
	if cookie, err := r.Cookie("session_id"); err == nil {
		// Cookie exists, verify if session is valid
		if user, err := h.sessions.GetUserFromSession(r.Context(), cookie.Value); err == nil && user != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	component := auth.LoginPage(auth.LoginProps{})
	component.Render(r.Context(), w)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember-me") == "on"

	// Validate input
	if len(password) < MinPasswordLength {
		auth.ErrorMessage("Password must be at least 6 characters long").Render(r.Context(), w)
		return
	}

	user, err := h.service.Login(r.Context(), email, password)
	if err != nil {
		var errMsg string
		if errors.Is(err, ErrInvalidCredentials) {
			errMsg = "Invalid email or password"
		} else {
			errMsg = "An error occurred while logging in"
		}
		auth.ErrorMessage(errMsg).Render(r.Context(), w)
		return
	}

	// Create session
	session, err := h.sessions.Create(r.Context(), user.ID, rememberMe)
	if err != nil {
		auth.ErrorMessage("Error creating session").Render(r.Context(), w)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.getSessionMaxAge(rememberMe),
	})

	// Redirect to dashboard
	w.Header().Set("HX-Redirect", "/dashboard")
}

func (h *Handler) getSessionMaxAge(rememberMe bool) int {
	if rememberMe {
		return 30 * 24 * 60 * 60 // 30 days
	}
	return 24 * 60 * 60 // 1 day
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

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember-me") == "on"

	// Validate input
	if len(password) < MinPasswordLength {
		auth.ErrorMessage("Password must be at least 6 characters long").Render(r.Context(), w)
		return
	}

	// Use the service's SignUp method
	if err := h.service.SignUp(r.Context(), email, password); err != nil {
		auth.ErrorMessage("Error creating account").Render(r.Context(), w)
		return
	}

	// Since SignUp doesn't return the user, we need to log them in
	user, err := h.service.Login(r.Context(), email, password)
	if err != nil {
		auth.ErrorMessage("Account created but could not log in").Render(r.Context(), w)
		return
	}

	// Create session
	session, err := h.sessions.Create(r.Context(), user.ID, rememberMe)
	if err != nil {
		auth.ErrorMessage("Account created but could not log in").Render(r.Context(), w)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.getSessionMaxAge(rememberMe),
	})

	if !strings.HasSuffix(email, "@merindol.co") {
		http.Redirect(w, r, "/auth/login?pending_verification=true", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/auth/login?signup_success=true", http.StatusSeeOther)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	// Delete session from store using the session ID from context
	if sessionID := h.sessions.GetSessionID(r.Context()); sessionID != "" {
		h.sessions.Invalidate(r.Context(), sessionID)
	}

	// Redirect to login page
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (h *Handler) handleVerify(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if err := h.service.VerifyEmail(r.Context(), token); err != nil {
		http.Error(w, "Invalid or expired verification link", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/auth/login?verified=true", http.StatusSeeOther)
}
