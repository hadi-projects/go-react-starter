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

type ArsipServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockArsipRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   ArsipService
}

func (s *ArsipServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockArsipRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewArsipService(s.mockRepo, s.mockCache)
}

func (s *ArsipServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ArsipServiceTestSuite) TestCreate_Success() {
	req := dto.CreateArsipRequest{
		Name: "test",
		Tanggal: "test",
	}

	entity := &entity.Arsip{
		Name: req.Name,
		Tanggal: req.Tanggal,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("arsip:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestArsipServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ArsipServiceTestSuite))
}
