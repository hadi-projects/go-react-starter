package service

import (
	"testing"

	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	mock_pkg "github.com/hadi-projects/go-react-starter/internal/mock/pkg"
	mock_repository "github.com/hadi-projects/go-react-starter/internal/mock/repository"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockUserRepo  *mock_repository.MockUserRepository
	mockTokenRepo *mock_repository.MockTokenRepository
	mockProducer  *mock_pkg.MockProducer
	mockMailer    *mock_pkg.MockMailer
	service       AuthService
	cfg           *config.Config
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockUserRepo = mock_repository.NewMockUserRepository(s.ctrl)
	s.mockTokenRepo = mock_repository.NewMockTokenRepository(s.ctrl)
	s.mockProducer = mock_pkg.NewMockProducer(s.ctrl)
	s.mockMailer = mock_pkg.NewMockMailer(s.ctrl)

	s.cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
		Kafka: config.KafkaConfig{
			Topic: "test-topic",
		},
	}

	logger.AuthLogger = zerolog.Nop()
	logger.SystemLogger = zerolog.Nop()
	s.service = NewAuthService(s.mockUserRepo, s.mockTokenRepo, s.mockProducer, s.mockMailer, s.cfg)
}

func (s *AuthServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *AuthServiceTestSuite) TestLogin_Success() {
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &entity.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role: entity.Role{
			Name: "admin",
		},
	}

	req := dto.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	s.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(user, nil)

	res, err := s.service.Login(req)
	s.Require().NoError(err)
	assert.NotNil(s.T(), res)
	assert.NotEmpty(s.T(), res.AccessToken)
}

func (s *AuthServiceTestSuite) TestLogin_InvalidPassword() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), 10)

	user := &entity.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	req := dto.LoginRequest{
		Email:    user.Email,
		Password: "wrong-password",
	}

	s.mockUserRepo.EXPECT().FindByEmail(req.Email).Return(user, nil)

	res, err := s.service.Login(req)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), "invalid email or password", err.Error())
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
