package service

import (
	"context"
	"gopkg.in/gomail.v2"
	"notification-service/config"
	"notification-service/internal/model"
)

type MailService struct {
	config config.Mail
	dialer *gomail.Dialer
}

func NewMailService(config config.Mail) *MailService {
	return &MailService{
		config: config,
		dialer: gomail.NewDialer(config.Host, config.Port, config.Sender, config.Password),
	}
}

func (m *MailService) Notify(ctx context.Context, message, title string, receiver model.Receiver) error {
	goMailMessage := gomail.NewMessage()
	goMailMessage.SetHeader("From", m.config.Sender)
	goMailMessage.SetHeader("To", receiver.Mail)
	goMailMessage.SetHeader("Subject", title)
	goMailMessage.SetBody("text", message)
	if err := m.dialer.DialAndSend(goMailMessage); err != nil {
		return err
	}

	return m.dialer.DialAndSend(goMailMessage)
}
