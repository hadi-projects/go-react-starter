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

type TestsajaServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockTestsajaRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   TestsajaService
}

func (s *TestsajaServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockTestsajaRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewTestsajaService(s.mockRepo, s.mockCache)
}

func (s *TestsajaServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *TestsajaServiceTestSuite) TestCreate_Success() {
	req := dto.CreateTestsajaRequest{
		Name: "test",
	}

	entity := &entity.Testsaja{
		Name: req.Name,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("testsaja:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestTestsajaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TestsajaServiceTestSuite))
}
