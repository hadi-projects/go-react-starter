package mailer

import (
	"context"
	"fmt"

	"time"

	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
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

func (m *mailer) SendEmail(ctx context.Context, to string, subject string, body string) error {
	start := time.Now()
	requestID := ""
	if rid, ok := ctx.Value(logger.CtxKeyRequestID).(string); ok {
		requestID = rid
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.Mail.FromAddress)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	err := m.dialer.DialAndSend(msg)
	elapsed := time.Since(start)

	status := 200
	if err != nil {
		status = 500
	}

	logger.SystemLogger.Info().
		Str("method", "SMTP:SEND").
		Str("path", to).
		Int("status_code", status).
		Int64("latency", elapsed.Milliseconds()).
		Str("request_body", subject).
		Msg("mailer operation")

	if logger.SystemLogRepo != nil {
		_ = logger.SystemLogRepo.Create(ctx, &logger.SystemLog{
			RequestID:   requestID,
			Method:      "SMTP:SEND",
			Path:        to,
			StatusCode:  status,
			Latency:     elapsed.Milliseconds(),
			RequestBody: subject,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
