package app

import (
	"AP2_assignment1/service/internal/notification/domain"
	"AP2_assignment1/service/internal/notification/infrastructure/email"
	"os"
)

type NotificationApp struct {
	EmailSender domain.EmailSender
}

func NewNotificationApp() *NotificationApp {
	var sender domain.EmailSender

	if os.Getenv("PROVIDER_MODE") == "REAL" {
		sender = &email.MockEmailSender{}
	} else {
		sender = &email.MockEmailSender{}
	}

	return &NotificationApp{
		EmailSender: sender,
	}
}
