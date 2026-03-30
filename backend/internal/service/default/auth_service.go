package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/kafka"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	Logout(ctx context.Context, req dto.LogoutRequest) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	producer  kafka.Producer
	mailer    mailer.Mailer
	config    *config.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	producer kafka.Producer,
	mailer mailer.Mailer,
	config *config.Config,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		producer:  producer,
		mailer:    mailer,
		config:    config,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Find user by email (simple, no preloads)
	user, err := s.userRepo.FindByEmailSimple(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 3. Fetch full user data including Role and Permissions for Token generation
	fullUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}
	user = fullUser

	// 4. Generate JWT Token
	var permissions []string
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Name)
	}

	claims := jwt.MapClaims{
		"sub":         user.ID,
		"email":       user.Email,
		"role":        user.Role.Name,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// Audit login
	logger.LogAudit(context.WithValue(context.WithValue(ctx, logger.CtxKeyUserID, user.ID), logger.CtxKeyUserEmail, user.Email), "LOGIN", "AUTH", fmt.Sprintf("%d", user.ID), "")

	// 5. Return response
	return &dto.LoginResponse{
		AccessToken:  signedToken,
		RefreshToken: "refresh-token-placeholder", // TODO: Implement refresh token
		User: dto.UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			RoleID:      user.RoleID,
			Role:        user.Role.Name,
			Permissions: permissions,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Return nil to avoid enumerating users
		return nil
	}

	// Audit forgot password
	logger.LogAudit(ctx, "FORGOT_PASSWORD", "AUTH", fmt.Sprintf("%d", user.ID), fmt.Sprintf("email: %s", req.Email))

	// 2. Generate token
	token := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	// 3. Save token
	resetToken := &entity.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := s.tokenRepo.Create(resetToken); err != nil {
		return err
	}

	// 4. Publish message to Kafka
	msg := map[string]string{
		"email": user.Email,
		"token": token,
	}

	// Use configured topic from config
	topic := s.config.Kafka.Topic
	if topic == "" {
		topic = "password-reset"
	}

	var publishErr error
	if s.producer != nil {
		publishErr = s.producer.Publish(topic, msg)
	} else {
		publishErr = errors.New("kafka producer is not initialized")
	}

	if publishErr != nil {
		logger.SystemLogger.Error().Err(publishErr).Msg("Failed to publish password reset message to Kafka. Falling back to direct email.")

		// Fallback: Send email via goroutine
		go func() {
			frontendURL := s.config.Frontend.URL
			if frontendURL == "" {
				frontendURL = "http://localhost:5173"
			}
			resetLink := frontendURL + "/reset-password?token=" + token
			body := mailer.GetResetPasswordEmailNative(resetLink)
			if err := s.mailer.SendEmail(user.Email, "Reset Password Request (Fallback)", body); err != nil {
				logger.SystemLogger.Error().Err(err).Str("email", user.Email).Msg("Failed to send fallback email")
			} else {
				logger.SystemLogger.Info().Str("email", user.Email).Msg("Fallback email sent successfully")
			}
		}()

		// Return nil to client as the request is accepted
		return nil
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	// 1. Find token
	resetToken, err := s.tokenRepo.FindByToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// 2. Check expiration
	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("invalid or expired token")
	}

	// 3. Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.Security.BCryptCost)
	if err != nil {
		return err
	}

	// 4. Update user password
	user := resetToken.User
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(&user); err != nil {
		return err
	}

	// 5. Delete token (and potentially all other tokens for this user)
	if err := s.tokenRepo.DeleteByUserID(user.ID); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Failed to delete reset tokens")
	}

	// Audit reset password
	logger.LogAudit(context.WithValue(context.WithValue(ctx, logger.CtxKeyUserID, user.ID), logger.CtxKeyUserEmail, user.Email), "RESET_PASSWORD", "AUTH", fmt.Sprintf("%d", user.ID), "")

	return nil
}

func (s *authService) Logout(ctx context.Context, req dto.LogoutRequest) error {
	userID, _ := ctx.Value(logger.CtxKeyUserID).(uint)

	// Audit logout
	logger.LogAudit(ctx, "LOGOUT", "AUTH", fmt.Sprintf("%d", userID), fmt.Sprintf("reason: %s", req.Reason))

	return nil
}
