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

type NewsServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockNewsRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   NewsService
}

func (s *NewsServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockNewsRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewNewsService(s.mockRepo, s.mockCache)
}

func (s *NewsServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *NewsServiceTestSuite) TestCreate_Success() {
	req := dto.CreateNewsRequest{
		Name: "test",
		Content: "test",
	}

	entity := &entity.News{
		Name: req.Name,
		Content: req.Content,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("news:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestNewsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NewsServiceTestSuite))
}
