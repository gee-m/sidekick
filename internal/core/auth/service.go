package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
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

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
	}
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	var user User
	err := s.db.QueryRow(ctx,
		"SELECT id, email, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("querying user: %w", err)
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

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, email, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating users: %w", err)
	}

	return users, nil
}
