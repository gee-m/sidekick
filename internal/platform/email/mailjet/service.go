package mailjet

import (
	"fmt"

	"appgents/internal/core/types"
	"context"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Config struct {
	APIKey    string
	SecretKey string
	FromEmail string
	FromName  string
}

type Service struct {
	client *mailjet.Client
	config Config
}

type Email struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

func NewService(config Config) *Service {
	return &Service{
		client: mailjet.NewMailjetClient(config.APIKey, config.SecretKey),
		config: config,
	}
}

func (s *Service) SendMail(email types.Email) error {
	message := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: s.config.FromEmail,
			Name:  s.config.FromName,
		},
		To: &mailjet.RecipientsV31{
			{
				Email: email.To,
			},
		},
		Subject:  email.Subject,
		TextPart: email.Text,
		HTMLPart: email.Text,
	}

	messages := mailjet.MessagesV31{Info: []mailjet.InfoMessagesV31{message}}
	_, err := s.client.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}

func (s *Service) SendVerificationEmail(ctx context.Context, to string, verificationURL string) error {
	email := types.Email{
		To:      to,
		Subject: "Verify your Appgents account",
		Text: fmt.Sprintf(`
<h3>Welcome to Appgents!</h3>
<p>Please click the link below to verify your email address:</p>
<p><a href="%s">Verify Email</a></p>
<p>If you didn't create an account, you can safely ignore this email.</p>
`, verificationURL),
	}

	return s.SendMail(email)
}

func (s *Service) SendPasswordResetEmail(ctx context.Context, to string, resetURL string) error {
	email := types.Email{
		To:      to,
		Subject: "Reset your Appgents password",
		Text: fmt.Sprintf(`
<h3>Password Reset Request</h3>
<p>Click the link below to reset your password:</p>
<p><a href="%s">Reset Password</a></p>
<p>If you didn't request a password reset, you can safely ignore this email.</p>
`, resetURL),
	}

	return s.SendMail(email)
}
