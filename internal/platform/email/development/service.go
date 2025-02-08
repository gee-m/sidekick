package development

import (
	"appgents/internal/core/types"
	"appgents/internal/platform/logging"
	"context"
)

type Service struct {
	logger *logging.Logger
}

func New(logger *logging.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) SendMail(email types.Email) error {
	s.logger.Info("Email would be sent", map[string]any{
		"to":      email.To,
		"subject": email.Subject,
		"body":    email.Text,
	})
	return nil
}

func (s *Service) SendVerificationEmail(ctx context.Context, to string, verificationURL string) error {
	s.logger.Info("Verification email would be sent", map[string]any{
		"to":              to,
		"verificationURL": verificationURL,
	})
	return nil
}

func (s *Service) SendPasswordResetEmail(ctx context.Context, to string, resetURL string) error {
	s.logger.Info("Password reset email would be sent", map[string]any{
		"to":       to,
		"resetURL": resetURL,
	})
	return nil
}
