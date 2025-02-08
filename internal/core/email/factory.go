package email

import (
	"appgents/internal/platform/email/development"
	"appgents/internal/platform/email/mailjet"
	"appgents/internal/platform/logging"
)

type Config struct {
	Environment string
	Mailjet     mailjet.Config
}

func NewEmailService(config Config, logger *logging.Logger) (Service, error) {
	// if config.Environment != "production" {
	return development.New(logger), nil
	// }

	// Production email service
	// return mailjet.New(config.Mailjet)
}
