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

type abcServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockAbcRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   AbcService
}

func (s *abcServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockAbcRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewAbcService(s.mockRepo, s.mockCache)
}

func (s *abcServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *abcServiceTestSuite) TestCreate_Success() {
	req := dto.CreateAbcRequest{
		Name: "test",
	}

	entity := &entity.Abc{
		Name: req.Name,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("abc:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestAbcServiceTestSuite(t *testing.T) {
	suite.Run(t, new(abcServiceTestSuite))
}
