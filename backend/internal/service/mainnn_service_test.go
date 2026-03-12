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

type MainnnServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockMainnnRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   MainnnService
}

func (s *MainnnServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockMainnnRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewMainnnService(s.mockRepo, s.mockCache)
}

func (s *MainnnServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *MainnnServiceTestSuite) TestCreate_Success() {
	req := dto.CreateMainnnRequest{
		Name: "test",
		Makananan: "test",
	}

	entity := &entity.Mainnn{
		Name: req.Name,
		Makananan: req.Makananan,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("mainnn:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestMainnnServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MainnnServiceTestSuite))
}
