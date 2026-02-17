package service

import (
	"errors"
	"testing"

	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	mock_cache "github.com/hadi-projects/go-react-starter/internal/mock/pkg/cache"
	mock_repository "github.com/hadi-projects/go-react-starter/internal/mock/repository"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockUserRepository
	mockCache *mock_cache.MockCacheService
	service   UserService
	cfg       *config.Config
}

func (s *UserServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockUserRepository(s.ctrl)
	s.mockCache = mock_cache.NewMockCacheService(s.ctrl)

	s.cfg = &config.Config{
		Security: config.SecurityConfig{
			BCryptCost: 10,
		},
	}

	// Initialize AuditLogger with Nop to avoid panics
	logger.AuditLogger = zerolog.Nop()

	s.service = NewUserService(s.mockRepo, s.cfg, s.mockCache)
}

func (s *UserServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *UserServiceTestSuite) TestRegister_Success() {
	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	// Expect check if email exists
	s.mockRepo.EXPECT().FindByEmail(req.Email).Return(nil, errors.New("not found"))

	// Expect check/get default role
	s.mockRepo.EXPECT().FindRoleByName("user").Return(&entity.Role{ID: 2, Name: "user"}, nil)

	// Expect create user
	s.mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *entity.User) error {
		user.ID = 1
		return nil
	})

	// Expect cache invalidation
	s.mockCache.EXPECT().DeletePattern("users:*").Return(nil)

	res, err := s.service.Register(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), uint(1), res.ID)
	assert.Equal(s.T(), req.Email, res.Email)
}

func (s *UserServiceTestSuite) TestRegister_EmailExists() {
	req := dto.RegisterRequest{
		Email: "existing@example.com",
	}

	s.mockRepo.EXPECT().FindByEmail(req.Email).Return(&entity.User{ID: 1}, nil)

	res, err := s.service.Register(req)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), "email already exists", err.Error())
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
