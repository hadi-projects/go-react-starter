package mailer

import (
	"fmt"

	"github.com/hadi-projects/go-react-starter/config"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendEmail(to string, subject string, body string) error
}

type mailer struct {
	dialer *gomail.Dialer
	cfg    *config.Config
}

func NewMailer(cfg *config.Config) Mailer {
	dialer := gomail.NewDialer(
		cfg.Mail.Host,
		cfg.Mail.Port,
		cfg.Mail.User,
		cfg.Mail.Password,
	)

	return &mailer{
		dialer: dialer,
		cfg:    cfg,
	}
}

func (m *mailer) SendEmail(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.Mail.FromAddress)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
