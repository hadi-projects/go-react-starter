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

type CookServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockCookRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   CookService
}

func (s *CookServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockCookRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewCookService(s.mockRepo, s.mockCache)
}

func (s *CookServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CookServiceTestSuite) TestCreate_Success() {
	req := dto.CreateCookRequest{
		Name: "test",
	}

	entity := &entity.Cook{
		Name: req.Name,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("cook:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestCookServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CookServiceTestSuite))
}
