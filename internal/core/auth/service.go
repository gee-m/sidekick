package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnverifiedEmail    = errors.New("email not verified")
	ErrInvalidToken       = errors.New("invalid verification token")
)

type Service struct {
	db *pgxpool.Pool
}

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type UserInfo struct {
	Email           string
	CreatedAt       time.Time
	VerificationURL string // Changed from pointer to string
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
	}
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	var (
		user       User
		verifiedAt *time.Time
	)
	err := s.db.QueryRow(ctx,
		"SELECT id, email, password_hash, verified_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &verifiedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("querying user: %w", err)
	}

	if !strings.HasSuffix(email, "@merindol.co") && verifiedAt == nil {
		return nil, ErrUnverifiedEmail
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// HashPassword is a helper function to hash passwords before storage
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hashedBytes), nil
}

func (s *Service) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		email,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("checking email existence: %w", err)
	}

	return exists, nil
}

func (s *Service) CreateUser(ctx context.Context, email, hashedPassword string) (*User, error) {
	var user User
	err := s.db.QueryRow(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email",
		email, hashedPassword,
	).Scan(&user.ID, &user.Email)

	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return &user, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]UserInfo, error) {
	query := `
        SELECT email, created_at, verification_token
        FROM users
        ORDER BY created_at DESC
    `
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []UserInfo
	for rows.Next() {
		var (
			user              UserInfo
			verificationToken *uuid.UUID
		)
		if err := rows.Scan(&user.Email, &user.CreatedAt, &verificationToken); err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}

		if verificationToken != nil {
			user.VerificationURL = fmt.Sprintf("/auth/verify/%s", verificationToken.String())
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Service) SignUp(ctx context.Context, email, password string) error {
	// Check if user already exists
	exists, err := s.EmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	needsVerification := !strings.HasSuffix(email, "@merindol.co")
	var verificationToken *uuid.UUID

	if needsVerification {
		token := uuid.New()
		verificationToken = &token
	}

	// Create user with verification token if needed
	query := `
        INSERT INTO users (email, password_hash, verification_token)
        VALUES ($1, $2, $3)
    `
	if _, err := s.db.Exec(ctx, query, email, hashedPassword, verificationToken.String()); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	query := `
        UPDATE users
        SET verified_at = NOW(), verification_token = NULL
        WHERE verification_token = $1
        RETURNING id
    `
	var id string
	if err := s.db.QueryRow(ctx, query, token).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrInvalidToken
		}
		return fmt.Errorf("verifying email: %w", err)
	}
	return nil
}
