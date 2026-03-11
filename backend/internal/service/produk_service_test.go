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

type ProdukServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockProdukRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   ProdukService
}

func (s *ProdukServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockProdukRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewProdukService(s.mockRepo, s.mockCache)
}

func (s *ProdukServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ProdukServiceTestSuite) TestCreate_Success() {
	req := dto.CreateProdukRequest{
		Name: "test",
		Harga: 5,
	}

	entity := &entity.Produk{
		Name: req.Name,
		Harga: req.Harga,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("produk:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestProdukServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProdukServiceTestSuite))
}
