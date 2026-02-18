package worker

import (
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

	logger.SystemLogger.Info().Str("email", data.Email).Msg("Processing reset password email")

	// Construct reset link (hardcoded for now, should come from config)
	// Assuming frontend is at localhost:3000
	// TODO: move base URL to config
	resetLink := "http://localhost:3000/reset-password?token=" + data.Token

	body := mailer.GetResetPasswordEmailNative(resetLink)

	if err := mailService.SendEmail(data.Email, "Reset Password Request", body); err != nil {
		logger.SystemLogger.Error().Err(err).Str("email", data.Email).Msg("Failed to send reset password email")
		return err
	}

	logger.SystemLogger.Info().Str("email", data.Email).Msg("Reset password email sent successfully")
	return nil
}
