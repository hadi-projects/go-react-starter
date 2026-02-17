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

type RoleServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockRoleRepository
	mockCache *mock_cache.MockCacheService
	service   RoleService
}

func (s *RoleServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockRoleRepository(s.ctrl)
	s.mockCache = mock_cache.NewMockCacheService(s.ctrl)

	logger.AuditLogger = zerolog.Nop()
	s.service = NewRoleService(s.mockRepo, s.mockCache)
}

func (s *RoleServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *RoleServiceTestSuite) TestCreate_Success() {
	req := dto.CreateRoleRequest{
		Name:          "admin",
		PermissionIDs: []uint{1},
	}

	s.mockRepo.EXPECT().Create(gomock.Any(), req.PermissionIDs).DoAndReturn(func(role *entity.Role, ids []uint) error {
		role.ID = 1
		role.Name = req.Name
		return nil
	})

	s.mockCache.EXPECT().DeletePattern("roles:*").Return(nil)
	s.mockRepo.EXPECT().FindByID(uint(1)).Return(&entity.Role{ID: 1, Name: "admin"}, nil)

	res, err := s.service.Create(req)
	s.Require().NoError(err)
	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), "admin", res.Name)
}

func TestRoleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}
