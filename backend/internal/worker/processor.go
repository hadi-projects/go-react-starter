package worker

import (
	"context"
	"encoding/json"

	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
)

type ResetPasswordPayload struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	AppName string `json:"app_name"`
	LogoURL string `json:"logo_url"`
}

func ProcessResetPassword(payload []byte, cfg *config.Config, mailService mailer.Mailer) error {
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

	// Construct reset link using configuration
	frontendURL := cfg.Frontend.URL
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	resetLink := frontendURL + "/reset-password?token=" + data.Token
	
	// Debug log to verify the generated link
	logger.SystemLogger.Info().
		Str("reset_link", resetLink).
		Str("app_name", data.AppName).
		Msg("Generated Reset Password Link")
	
	// The processor currently doesn't have access to SettingService
	// We should pass the logo URL in the message payload or fetch it here
	// For now, we'll use an empty logo URL or pass it from the producer
	body := mailer.GetResetPasswordEmailNative(resetLink, data.AppName, data.LogoURL)

	if err := mailService.SendEmail(context.Background(), data.Email, "Reset Password Request", body); err != nil {
		return err
	}

	return nil
}
