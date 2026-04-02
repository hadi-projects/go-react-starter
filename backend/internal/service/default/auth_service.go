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
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/kafka"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	Logout(ctx context.Context, req dto.LogoutRequest) error
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	Enroll2FA(ctx context.Context, userID uint) (*dto.TwoFAEnrollResponse, error)
	Confirm2FA(ctx context.Context, userID uint, req dto.TwoFAConfirmRequest) error
	Disable2FA(ctx context.Context, userID uint, req dto.TwoFADisableRequest) error
	Verify2FA(ctx context.Context, req dto.TwoFAVerifyRequest) (*dto.LoginResponse, error)
	Request2FAReset(ctx context.Context, req dto.TwoFAResetRequest) error
	Confirm2FAReset(ctx context.Context, req dto.TwoFAResetConfirmRequest) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	producer  kafka.Producer
	mailer    mailer.Mailer
	config    *config.Config
	cache     cache.CacheService
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	producer kafka.Producer,
	mailer mailer.Mailer,
	config *config.Config,
	cache cache.CacheService,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		producer:  producer,
		mailer:    mailer,
		config:    config,
		cache:     cache,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Find user by email (simple, no preloads)
	user, err := s.userRepo.FindByEmailSimple(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status == "freezed" {
		return nil, errors.New("your account is frozen, please contact administrator")
	}
	if user.Status == "pending" {
		return nil, errors.New("your account is pending approval")
	}
	// 2. Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 3. Fetch full user data including Role and Permissions for Token generation
	fullUser, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user = fullUser

	// 4. Generate JWT Tokens
	// 4a. If 2FA is enabled, issue a temp token instead of JWT
	if user.TwoFAEnabled {
		tempToken := uuid.New().String()
		tempKey := "2fa_temp:" + tempToken
		s.cache.Set(ctx, tempKey, fmt.Sprintf("%d", user.ID), 5*time.Minute)
		return &dto.LoginResponse{
			Requires2FA: true,
			TempToken:   tempToken,
		}, nil
	}

	// 4b. No 2FA: Generate JWT Tokens directly
	var permissionsMask uint64
	for _, p := range user.Role.Permissions {
		if p.ID <= 64 {
			permissionsMask |= (uint64(1) << (p.ID - 1))
		}
	}

	accessToken, err := s.generateAccessToken(user, permissionsMask)
	if err != nil {
		return nil, err
	}

	var refreshTokenStr string
	if req.RememberMe {
		refreshTokenStr = uuid.New().String()
		expirationDays := 7
		if s.config.JWT.RefreshExpirationTime != "" {
			fmt.Sscanf(s.config.JWT.RefreshExpirationTime, "%dh", &expirationDays)
		}
		expiresAt := time.Now().Add(time.Hour * time.Duration(24*expirationDays))
		rt := &entity.RefreshToken{
			UserID:    user.ID,
			Token:     refreshTokenStr,
			ExpiresAt: expiresAt,
		}
		if err := s.tokenRepo.CreateRefreshToken(ctx, rt); err != nil {
			return nil, err
		}
	}

	logger.LogAudit(context.WithValue(context.WithValue(ctx, logger.CtxKeyUserID, user.ID), logger.CtxKeyUserEmail, user.Email), "LOGIN", "AUTH", fmt.Sprintf("%d", user.ID), fmt.Sprintf("RememberMe: %v", req.RememberMe))

	// 5. Return response
	userResp := &dto.AuthUserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		RoleID:          user.RoleID,
		Role:            user.Role.Name,
		PermissionsMask: permissionsMask,
		Status:          user.Status,
		TwoFAEnabled:    user.TwoFAEnabled,
	}
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         userResp,
	}, nil
}

func (s *authService) generateAccessToken(user *entity.User, permissionsMask uint64) (string, error) {
	claims := jwt.MapClaims{
		"sub":              user.ID,
		"email":            user.Email,
		"role":             user.Role.Name,
		"permissions_mask": permissionsMask,
		"exp":              time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// 1. Find refresh token
	rt, err := s.tokenRepo.FindByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 2. Check expiration
	if time.Now().After(rt.ExpiresAt) {
		s.tokenRepo.DeleteRefreshToken(ctx, req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	if rt.User.Status == "freezed" {
		return nil, errors.New("your account is frozen, please contact administrator")
	}
	if rt.User.Status == "pending" {
		return nil, errors.New("your account is pending approval")
	}

	// 3. Generate new Access Token
	var permissionsMask uint64
	for _, p := range rt.User.Role.Permissions {
		if p.ID <= 64 {
			permissionsMask |= (uint64(1) << (p.ID - 1))
		}
	}

	accessToken, err := s.generateAccessToken(&rt.User, permissionsMask)
	if err != nil {
		return nil, err
	}

	// 4. Return response
	userResp := &dto.AuthUserResponse{
		ID:              rt.User.ID,
		Name:            rt.User.Name,
		Email:           rt.User.Email,
		RoleID:          rt.User.RoleID,
		Role:            rt.User.Role.Name,
		PermissionsMask: permissionsMask,
		Status:          rt.User.Status,
		TwoFAEnabled:    rt.User.TwoFAEnabled,
	}
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rt.Token,
		User:         userResp,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
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

	if err := s.tokenRepo.Create(ctx, resetToken); err != nil {
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
		publishErr = s.producer.Publish(ctx, topic, msg)
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
			if err := s.mailer.SendEmail(context.Background(), user.Email, "Reset Password Request (Fallback)", body); err != nil {
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
	resetToken, err := s.tokenRepo.FindByToken(ctx, req.Token)
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
	if err := s.userRepo.Update(ctx, &user); err != nil {
		return err
	}

	// 5. Delete token (and potentially all other tokens for this user)
	if err := s.tokenRepo.DeleteByUserID(ctx, user.ID); err != nil {
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

// Enroll2FA generates a new TOTP secret for the user (does not activate 2FA yet)
func (s *authService) Enroll2FA(ctx context.Context, userID uint) (*dto.TwoFAEnrollResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.TwoFAEnabled {
		return nil, errors.New("2FA is already enabled")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GoReactApp",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	// Save the secret (not enabled yet, needs confirmation)
	user.TwoFASecret = key.Secret()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.TwoFAEnrollResponse{
		Secret: key.Secret(),
		QRURL:  key.URL(),
	}, nil
}

// Confirm2FA activates 2FA after verifying the first code
func (s *authService) Confirm2FA(ctx context.Context, userID uint, req dto.TwoFAConfirmRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.TwoFAEnabled {
		return errors.New("2FA is already enabled")
	}
	if user.TwoFASecret == "" {
		return errors.New("no 2FA secret found, please enroll first")
	}

	// Validate TOTP code
	valid := totp.Validate(req.Code, user.TwoFASecret)
	if !valid {
		return errors.New("invalid 2FA code")
	}

	user.TwoFAEnabled = true
	return s.userRepo.Update(ctx, user)
}

// Disable2FA verifies the code then turns off 2FA
func (s *authService) Disable2FA(ctx context.Context, userID uint, req dto.TwoFADisableRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !user.TwoFAEnabled {
		return errors.New("2FA is not enabled")
	}

	// Validate TOTP code
	valid := totp.Validate(req.Code, user.TwoFASecret)
	if !valid {
		return errors.New("invalid 2FA code")
	}

	user.TwoFAEnabled = false
	user.TwoFASecret = ""
	user.TwoFACounter = 0
	return s.userRepo.Update(ctx, user)
}

// Verify2FA exchanges a temp_token + code for a real JWT
func (s *authService) Verify2FA(ctx context.Context, req dto.TwoFAVerifyRequest) (*dto.LoginResponse, error) {
	tempKey := "2fa_temp:" + req.TempToken
	var userIDStr string
	if err := s.cache.Get(ctx, tempKey, &userIDStr); err != nil || userIDStr == "" {
		return nil, errors.New("invalid or expired 2FA session, please login again")
	}

	var userID uint
	fmt.Sscanf(userIDStr, "%d", &userID)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate TOTP code
	valid := totp.Validate(req.Code, user.TwoFASecret)
	if !valid {
		return nil, errors.New("invalid 2FA code")
	}

	// Delete the temp token
	s.cache.Delete(ctx, tempKey)

	// Generate the full JWT
	return s.generateLoginResponse(ctx, user)
}

// Helper to generate standard login response with audit logging
func (s *authService) generateLoginResponse(ctx context.Context, user *entity.User) (*dto.LoginResponse, error) {
	var permissionsMask uint64
	for _, p := range user.Role.Permissions {
		if p.ID <= 64 {
			permissionsMask |= (uint64(1) << (p.ID - 1))
		}
	}

	accessToken, err := s.generateAccessToken(user, permissionsMask)
	if err != nil {
		return nil, err
	}

	logger.LogAudit(context.WithValue(ctx, logger.CtxKeyUserID, user.ID), "LOGIN_2FA", "AUTH", fmt.Sprintf("%d", user.ID), "")

	userResp := &dto.AuthUserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		RoleID:          user.RoleID,
		Role:            user.Role.Name,
		PermissionsMask: permissionsMask,
		Status:          user.Status,
		TwoFAEnabled:    user.TwoFAEnabled,
	}
	return &dto.LoginResponse{
		AccessToken: accessToken,
		User:        userResp,
	}, nil
}

func (s *authService) Request2FAReset(ctx context.Context, req dto.TwoFAResetRequest) error {
	tempKey := "2fa_temp:" + req.TempToken
	var userIDStr string
	if err := s.cache.Get(ctx, tempKey, &userIDStr); err != nil || userIDStr == "" {
		return errors.New("invalid or expired 2FA session")
	}

	var userID uint
	fmt.Sscanf(userIDStr, "%d", &userID)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Delete any existing reset tokens
	s.tokenRepo.DeleteTwoFAResetTokenByUserID(ctx, user.ID)

	token := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	resetToken := &entity.TwoFAResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := s.tokenRepo.CreateTwoFAResetToken(ctx, resetToken); err != nil {
		return err
	}

	frontendURL := s.config.Frontend.URL
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	resetLink := frontendURL + "/twofa/reset-confirm?token=" + token

	go func() {
		body := mailer.GetTwoFAResetEmailNative(resetLink)
		if err := s.mailer.SendEmail(context.Background(), user.Email, "Reset Two-Factor Authentication", body); err != nil {
			logger.SystemLogger.Error().Err(err).Str("email", user.Email).Msg("Failed to send 2FA reset email")
		} else {
			logger.SystemLogger.Info().Str("email", user.Email).Msg("2FA reset email sent successfully")
		}
	}()

	return nil
}

func (s *authService) Confirm2FAReset(ctx context.Context, req dto.TwoFAResetConfirmRequest) error {
	resetToken, err := s.tokenRepo.FindByTwoFAResetToken(ctx, req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("invalid or expired token")
	}

	user := resetToken.User
	user.TwoFAEnabled = false
	user.TwoFASecret = ""
	user.TwoFACounter = 0

	if err := s.userRepo.Update(ctx, &user); err != nil {
		return err
	}

	if err := s.tokenRepo.DeleteTwoFAResetTokenByUserID(ctx, user.ID); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Failed to delete 2FA reset tokens")
	}

	logger.LogAudit(context.WithValue(context.WithValue(ctx, logger.CtxKeyUserID, user.ID), logger.CtxKeyUserEmail, user.Email), "RESET_2FA", "AUTH", fmt.Sprintf("%d", user.ID), "")

	return nil
}
