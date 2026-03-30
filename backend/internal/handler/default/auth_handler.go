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
	ForgotPassword(c *gin.Context)
	ResetPassword(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
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
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Login failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Login(c.Request.Context(), loginReq)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Login failed: service error")
		response.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

func (h *authHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), req); err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("ForgotPassword failed")
		// Always return success to avoid leaking internal errors or user existence
		response.Success(c, http.StatusOK, "If your email is registered, you will receive a password reset link.", nil)
		return
	}

	response.Success(c, http.StatusOK, "If your email is registered, you will receive a password reset link.", nil)
}

func (h *authHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req); err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("ResetPassword failed")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *authHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	// We use ShouldBindJSON but don't strictly require it for logout
	_ = c.ShouldBindJSON(&req)

	if err := h.service.Logout(c.Request.Context(), req); err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Logout failed")
		response.Error(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	response.Success(c, http.StatusOK, "Logout successful", nil)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed", res)
}
