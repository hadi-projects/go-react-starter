package handler

import (
	"bytes"
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

type UserHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockService *mock_service.MockUserService
	handler     UserHandler
}

func (s *UserHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockService = mock_service.NewMockUserService(s.ctrl)
	s.handler = NewUserHandler(s.mockService)

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize loggers to avoid panics
	logger.SystemLogger = zerolog.Nop()
}

func (s *UserHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *UserHandlerTestSuite) TestRegister_Success() {
	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}
	body, _ := json.Marshal(req)

	res := &dto.UserResponse{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
	}

	s.mockService.EXPECT().Register(gomock.Any(), req).Return(res, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	s.handler.Register(c)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	meta := resBody["meta"].(map[string]interface{})
	assert.Equal(s.T(), "success", meta["status"])
	assert.Equal(s.T(), "User registered successfully", meta["message"])
}

func (s *UserHandlerTestSuite) TestGetAll_Success() {
	pagination := &dto.PaginationResponse{
		Data: []dto.UserResponse{
			{ID: 1, Name: "User 1"},
		},
		Meta: dto.PaginationMeta{
			CurrentPage: 1,
			TotalPages:  1,
			TotalData:   1,
			Limit:       10,
		},
	}

	s.mockService.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(pagination, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users?page=1&limit=10", nil)

	s.handler.GetAll(c)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)
	meta := resBody["meta"].(map[string]interface{})
	assert.Equal(s.T(), "success", meta["status"])
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
