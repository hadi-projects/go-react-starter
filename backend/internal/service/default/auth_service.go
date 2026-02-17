package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{userRepo: userRepo, config: config}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		logger.AuthLogger.Warn().
			Str("email", req.Email).
			Str("action", "login").
			Msg("login failed: user not found")
		return nil, errors.New("invalid email or password")
	}

	// 2. Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.AuthLogger.Warn().
			Uint("user_id", user.ID).
			Str("email", user.Email).
			Str("action", "login").
			Msg("login failed: invalid password")
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT Token
	var permissions []string
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Name)
	}

	claims := jwt.MapClaims{
		"sub":         user.ID,
		"role":        user.Role.Name,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, err
	}

	logger.AuthLogger.Info().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Str("action", "login").
		Msg("login successful")

	// 4. Return response
	return &dto.LoginResponse{
		AccessToken:  signedToken,
		RefreshToken: "refresh-token-placeholder", // TODO: Implement refresh token
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			RoleID:    user.RoleID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
