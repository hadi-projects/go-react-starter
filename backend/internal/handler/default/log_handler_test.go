package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	mock_service "github.com/hadi-projects/go-react-starter/internal/mock/service"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type LogHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockService *mock_service.MockLogService
	handler     LogHandler
}

func (s *LogHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockService = mock_service.NewMockLogService(s.ctrl)
	s.handler = NewLogHandler(s.mockService)

	gin.SetMode(gin.TestMode)
	logger.SystemLogger = zerolog.Nop()
}

func (s *LogHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *LogHandlerTestSuite) TestGetLogs_AdminSuccess() {
	// Setup context with admin permissions and user_id
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/logs?type=audit", nil)
	c.Set("user_id", uint(1))
	c.Set("permissions", []string{"get-all-logs"})

	// Mock GetLogs
	expectedLogs := []dto.LogResponse{
		{Message: "test log"},
	}
	s.mockService.EXPECT().GetLogs(gomock.Any()).Return(expectedLogs, nil)

	s.handler.GetLogs(c)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	meta := resBody["meta"].(map[string]interface{})
	assert.Equal(s.T(), "success", meta["status"])
}

func TestLogHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(LogHandlerTestSuite))
}
