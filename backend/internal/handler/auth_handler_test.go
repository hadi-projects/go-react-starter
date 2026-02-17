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

type AuthHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockService *mock_service.MockAuthService
	handler     AuthHandler
}

func (s *AuthHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockService = mock_service.NewMockAuthService(s.ctrl)
	s.handler = NewAuthHandler(s.mockService)

	gin.SetMode(gin.TestMode)
	logger.SystemLogger = zerolog.Nop()
}

func (s *AuthHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *AuthHandlerTestSuite) TestLogin_Success() {
	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(req)

	res := &dto.LoginResponse{
		AccessToken: "test-token",
		User: dto.UserResponse{
			ID:    1,
			Email: req.Email,
		},
	}

	s.mockService.EXPECT().Login(req).Return(res, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	s.handler.Login(c)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	meta := resBody["meta"].(map[string]interface{})
	assert.Equal(s.T(), "success", meta["status"])
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
