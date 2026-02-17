package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type AuthHandler interface {
	Login(c *gin.Context)
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{service: service}
}

func (h *authHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Login failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Login(loginReq)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Login failed: service error")
		response.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}
