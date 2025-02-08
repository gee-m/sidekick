package auth

import (
	"net/http"

	authn "appgents/internal/core/auth"
	"appgents/internal/platform/logging"
	"appgents/web/templates/auth"
)

type Handler struct {
	logger      *logging.Logger
	authService *authn.Service
}

func New(logger *logging.Logger, authService *authn.Service) *Handler {
	return &Handler{
		logger:      logger,
		authService: authService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleLoginPage(w, r)
	case http.MethodPost:
		h.handleLogin(w, r)
	default:
		h.logger.Warn("method_not_allowed", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("serving_login_page", map[string]any{
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
	})

	err := auth.Login().Render(r.Context(), w)
	if err != nil {
		h.logger.Error("failed_to_render_login_page", err, map[string]any{
			"remote_addr": r.RemoteAddr,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("login_attempt", map[string]any{
		"remote_addr": r.RemoteAddr,
		"email":       r.FormValue("email"), // Note: In production, be careful logging PII
	})

	// For now, just return an error message
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Login functionality not implemented yet"))

	h.logger.Info("login_failed", map[string]any{
		"reason":      "not_implemented",
		"remote_addr": r.RemoteAddr,
	})
}
