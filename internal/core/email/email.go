package email

import (
	"appgents/internal/core/types"
	"fmt"
)

// Service defines the email service interface
type Service interface {
	SendMail(email types.Email) error
}

// Sender handles sending domain-specific emails
type Sender struct {
	service Service
	baseURL string
}

func NewSender(service Service, baseURL string) *Sender {
	return &Sender{
		service: service,
		baseURL: baseURL,
	}
}

func (s *Sender) SendVerification(to, token string) error {
	verifyURL := fmt.Sprintf("%s/auth/verify?token=%s", s.baseURL, token)

	email := types.Email{
		To:      to,
		Subject: "Verify your Appgents account",
		Text: fmt.Sprintf(`
<h3>Welcome to Appgents!</h3>
<p>Please click the link below to verify your email address:</p>
<p><a href="%s">Verify Email</a></p>
<p>If you didn't create an account, you can safely ignore this email.</p>
`, verifyURL),
	}

	return s.service.SendMail(email)
}
