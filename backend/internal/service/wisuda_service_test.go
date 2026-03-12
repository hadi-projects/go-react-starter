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

type WisudaServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockWisudaRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   WisudaService
}

func (s *WisudaServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockWisudaRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewWisudaService(s.mockRepo, s.mockCache)
}

func (s *WisudaServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *WisudaServiceTestSuite) TestCreate_Success() {
	req := dto.CreateWisudaRequest{
		Name: "test",
	}

	entity := &entity.Wisuda{
		Name: req.Name,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("wisuda:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestWisudaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WisudaServiceTestSuite))
}
