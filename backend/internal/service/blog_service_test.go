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

type BlogServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockBlogRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   BlogService
}

func (s *BlogServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockBlogRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewBlogService(s.mockRepo, s.mockCache)
}

func (s *BlogServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *BlogServiceTestSuite) TestCreate_Success() {
	req := dto.CreateBlogRequest{
		Name: "test",
		Content: "test",
		Thumbnail: "test",
	}

	entity := &entity.Blog{
		Name: req.Name,
		Content: req.Content,
		Thumbnail: req.Thumbnail,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("blog:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestBlogServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BlogServiceTestSuite))
}
