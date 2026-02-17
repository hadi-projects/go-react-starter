package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mock_cache "github.com/hadi-projects/go-react-starter/internal/mock/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CacheHandlerTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockCache *mock_cache.MockCacheService
	handler   CacheHandler
}

func (s *CacheHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockCache = mock_cache.NewMockCacheService(s.ctrl)
	s.handler = NewCacheHandler(s.mockCache)

	gin.SetMode(gin.TestMode)
	logger.SystemLogger = zerolog.Nop()
}

func (s *CacheHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CacheHandlerTestSuite) TestClearAll_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	s.mockCache.EXPECT().FlushAll().Return(nil)

	s.handler.ClearAll(c)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	meta := resBody["meta"].(map[string]interface{})
	assert.Equal(s.T(), "success", meta["status"])
}

func TestCacheHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CacheHandlerTestSuite))
}
