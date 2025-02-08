package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"appgents/internal/core/email"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type Service struct {
	userRepo     UserRepository
	emailSender  *email.Sender
	tokenTimeout time.Duration
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	UpdateVerified(userID string) error
	FindByVerificationToken(token string) (*User, error)
}

type User struct {
	ID                string
	Email             string
	PasswordHash      string
	Verified          bool
	VerificationToken string
	TokenExpiry       time.Time
}

func NewService(userRepo UserRepository, emailSender *email.Sender) *Service {
	return &Service{
		userRepo:     userRepo,
		emailSender:  emailSender,
		tokenTimeout: 24 * time.Hour,
	}
}

func (s *Service) Register(email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if email exists
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return ErrEmailTaken
	}

	// Generate verification token
	token, err := generateToken()
	if err != nil {
		return err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &User{
		Email:             email,
		PasswordHash:      string(hash),
		Verified:          false,
		VerificationToken: token,
		TokenExpiry:       time.Now().Add(s.tokenTimeout),
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// Send verification email
	return s.emailSender.SendVerification(email, token)
}

func (s *Service) Authenticate(email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return ErrInvalidCredentials
	}

	if !user.Verified {
		return ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

func (s *Service) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	if time.Now().After(user.TokenExpiry) {
		return ErrInvalidToken
	}

	return s.userRepo.UpdateVerified(user.ID)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
