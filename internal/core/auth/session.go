package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type contextKey string

const (
	UserContextKey      contextKey = "user"
	SessionIDContextKey contextKey = "session_id"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionManager struct {
	db *pgxpool.Pool
}

func NewSessionManager(db *pgxpool.Pool) *SessionManager {
	return &SessionManager{db: db}
}

func (sm *SessionManager) Create(ctx context.Context, userID string, rememberMe bool) (*Session, error) {
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		ExpiresAt: sm.getExpiryTime(rememberMe),
	}

	_, err := sm.db.Exec(ctx,
		"INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)",
		session.ID, session.UserID, session.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return session, nil
}

func (sm *SessionManager) getExpiryTime(rememberMe bool) time.Time {
	if rememberMe {
		return time.Now().AddDate(0, 0, 30) // 30 days
	}
	return time.Now().AddDate(0, 0, 1) // 1 day
}

func (sm *SessionManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		user, err := sm.GetUserFromSession(r.Context(), cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// Add both user and session ID to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, SessionIDContextKey, cookie.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sm *SessionManager) GetUserFromSession(ctx context.Context, sessionID string) (*User, error) {
	var user User
	err := sm.db.QueryRow(ctx, `
		SELECT u.id, u.email, u.created_at
		FROM users u
		JOIN sessions s ON s.user_id = u.id
		WHERE s.id = $1 AND s.expires_at > NOW()
	`, sessionID).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("getting user from session: %w", err)
	}

	return &user, nil
}

func (sm *SessionManager) Invalidate(ctx context.Context, sessionID string) error {
	_, err := sm.db.Exec(ctx, "DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return fmt.Errorf("invalidating session: %w", err)
	}
	return nil
}

func (sm *SessionManager) GetUser(ctx context.Context) *User {
	val := ctx.Value(UserContextKey)
	if val == nil {
		return nil
	}
	user, ok := val.(*User)
	if !ok {
		return nil
	}
	return user
}

func (sm *SessionManager) GetSessionID(ctx context.Context) string {
	val := ctx.Value(SessionIDContextKey)
	if val == nil {
		return ""
	}
	sessionID, ok := val.(string)
	if !ok {
		return ""
	}
	return sessionID
}
