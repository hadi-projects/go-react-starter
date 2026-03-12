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

type SupplierServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockRepo  *mock_repository.MockSupplierRepository
	mockCache *mock_pkg_cache.MockCacheService
	service   SupplierService
}

func (s *SupplierServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mock_repository.NewMockSupplierRepository(s.ctrl)
	s.mockCache = mock_pkg_cache.NewMockCacheService(s.ctrl)
	s.service = NewSupplierService(s.mockRepo, s.mockCache)
}

func (s *SupplierServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SupplierServiceTestSuite) TestCreate_Success() {
	req := dto.CreateSupplierRequest{
		Name: "test",
		Contact: "test",
		Address: "test",
	}

	entity := &entity.Supplier{
		Name: req.Name,
		Contact: req.Contact,
		Address: req.Address,
	}

	s.mockRepo.EXPECT().Create(entity).Return(nil)
	s.mockCache.EXPECT().DeletePattern("suppliers:*").Return(nil)

	res, err := s.service.Create(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestSupplierServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SupplierServiceTestSuite))
}
