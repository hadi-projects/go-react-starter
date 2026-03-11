package service

import (
	"testing"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	mock_repository "github.com/hadi-projects/go-react-starter/internal/mock/repository"
	mock_pkg_cache "github.com/hadi-projects/go-react-starter/internal/mock/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AdminServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockAdminRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   AdminService
}

func (s *AdminServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockAdminRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewAdminService(s.mockRepo, s.mockCache)
}

func (s *AdminServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *AdminServiceTestSuite) TestCreate_Success() {
	req := dto.CreateAdminRequest{
		Name: "test",
		Email: "test",
	}

	entity := &entity.Admin{
		Name: req.Name,
		Email: req.Email,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("admin:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestAdminServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AdminServiceTestSuite))
}
