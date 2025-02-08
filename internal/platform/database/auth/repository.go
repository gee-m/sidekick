package auth

import (
	"database/sql"
	"errors"

	"appgents/internal/core/auth"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByEmail(email string) (*auth.User, error) {
	var user auth.User
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, verified, verification_token, token_expiry
		FROM users
		WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Verified,
		&user.VerificationToken,
		&user.TokenExpiry,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(user *auth.User) error {
	return r.db.QueryRow(`
		INSERT INTO users (
			email,
			password_hash,
			verified,
			verification_token,
			token_expiry
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		user.Email,
		user.PasswordHash,
		user.Verified,
		user.VerificationToken,
		user.TokenExpiry,
	).Scan(&user.ID)
}

func (r *Repository) UpdateVerified(userID string) error {
	result, err := r.db.Exec(`
		UPDATE users
		SET verified = true,
			verification_token = NULL,
			token_expiry = NULL
		WHERE id = $1`,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return auth.ErrInvalidToken
	}
	return nil
}

func (r *Repository) FindByVerificationToken(token string) (*auth.User, error) {
	var user auth.User
	err := r.db.QueryRow(`
		SELECT id, email, verified, token_expiry
		FROM users
		WHERE verification_token = $1`,
		token,
	).Scan(&user.ID, &user.Email, &user.Verified, &user.TokenExpiry)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrInvalidToken
		}
		return nil, err
	}

	return &user, nil
}
