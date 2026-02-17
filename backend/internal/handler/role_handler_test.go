package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	mock_service "github.com/hadi-projects/go-react-starter/internal/mock/service"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RoleHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockService *mock_service.MockRoleService
	handler     RoleHandler
}

func (s *RoleHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockService = mock_service.NewMockRoleService(s.ctrl)
	s.handler = NewRoleHandler(s.mockService)

	gin.SetMode(gin.TestMode)
	logger.SystemLogger = zerolog.Nop()
}

func (s *RoleHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *RoleHandlerTestSuite) TestCreate_Success() {
	req := dto.CreateRoleRequest{
		Name: "admin",
	}
	body, _ := json.Marshal(req)

	res := &dto.RoleResponse{
		ID:   1,
		Name: req.Name,
	}

	s.mockService.EXPECT().Create(req).Return(res, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/roles", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	s.handler.Create(c)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	meta := resBody["meta"].(map[string]interface{})
	assert.Equal(s.T(), "success", meta["status"])
}

func TestRoleHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(RoleHandlerTestSuite))
}
