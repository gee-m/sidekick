package dashboard

import (
	"net/http"
	"strings"

	"github.com/gee-m/sidekick/internal/core/auth"
	"github.com/gee-m/sidekick/web/templates/dashboard"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *auth.Service
}

func NewHandler(service *auth.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(requireMerindolEmail)
	r.Get("/", h.handleDashboard)
	return r
}

func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	props := dashboard.DashboardProps{
		Users: make([]dashboard.UserInfo, len(users)),
	}

	for i, user := range users {
		props.Users[i] = dashboard.UserInfo{
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}

	dashboard.DashboardPage(props).Render(r.Context(), w)
}

func requireMerindolEmail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(auth.UserContextKey).(*auth.User)
		if !strings.HasSuffix(user.Email, "@merindol.co") {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
