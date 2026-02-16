package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/service"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type UserHandler interface {
	Register(c *gin.Context)
	Me(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{service: service}
}

func (h *userHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Register failed: invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Register(req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Register failed: service error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": res,
	})
}

func (h *userHandler) Me(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		logger.SystemLogger.Error().Msg("Me failed: user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := val.(uint)
	if !ok {
		logger.SystemLogger.Error().Msg("Me failed: invalid user_id type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	res, err := h.service.GetMe(userID)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("user_id", userID).Msg("Me failed: user not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func (h *userHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("GetAll users failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *userHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Update user failed: invalid ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Update user failed: invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Update user failed: service error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *userHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Delete user failed: invalid ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Delete user failed: service error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
