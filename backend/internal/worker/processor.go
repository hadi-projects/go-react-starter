package worker

import (
	"context"
	"encoding/json"

	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
)

type ResetPasswordPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func ProcessResetPassword(payload []byte, mailService mailer.Mailer) error {
	var data ResetPasswordPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	// This is now redundant with mailer's structured logging, 
	// but we can add a high-level one for the worker task itself.
	logger.SystemLogger.Info().
		Str("method", "WORKER:RESET_PASSWORD").
		Str("path", data.Email).
		Int("status_code", 200).
		Str("request_body", string(payload)).
		Msg("worker operation")

	// Construct reset link
	resetLink := "http://localhost:3000/reset-password?token=" + data.Token
	body := mailer.GetResetPasswordEmailNative(resetLink)

	if err := mailService.SendEmail(context.Background(), data.Email, "Reset Password Request", body); err != nil {
		return err
	}

	return nil
}
