package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
	}
}

// Login authenticates a user
func (s *Service) Login(username, password string) error {
	// TODO: Implement actual login logic
	return nil
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
