package service

import (
	"testing"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	mock_pkg_cache "github.com/hadi-projects/go-react-starter/internal/mock/pkg/cache"
	mock_repository "github.com/hadi-projects/go-react-starter/internal/mock/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TestduaServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockTestduaRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   TestduaService
}

func (s *TestduaServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockTestduaRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewTestduaService(s.mockRepo, s.mockCache)
}

func (s *TestduaServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *TestduaServiceTestSuite) TestCreate_Success() {
	req := dto.CreateTestduaRequest{
		Name: "test",
	}

	entity := &entity.Testdua{
		Name: req.Name,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("testdua:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestTestduaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TestduaServiceTestSuite))
}
