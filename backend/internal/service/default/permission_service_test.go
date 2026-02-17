package service

import (
	"testing"

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

type PermissionServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockPermissionRepository
	mockCache *mock_cache.MockCacheService
	service   PermissionService
}

func (s *PermissionServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockPermissionRepository(s.ctrl)
	s.mockCache = mock_cache.NewMockCacheService(s.ctrl)

	logger.AuditLogger = zerolog.Nop()
	s.service = NewPermissionService(s.mockRepo, s.mockCache)
}

func (s *PermissionServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *PermissionServiceTestSuite) TestCreate_Success() {
	req := dto.CreatePermissionRequest{
		Name: "test-permission",
	}

	s.mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *entity.Permission) error {
		p.ID = 1
		return nil
	})
	s.mockCache.EXPECT().DeletePattern("permissions:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), "test-permission", res.Name)
}

func TestPermissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionServiceTestSuite))
}
