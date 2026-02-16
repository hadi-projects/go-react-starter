package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/service"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type RoleHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type roleHandler struct {
	service service.RoleService
}

func NewRoleHandler(service service.RoleService) RoleHandler {
	return &roleHandler{service: service}
}

func (h *roleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create role failed: invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create role failed: service error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *roleHandler) GetAll(c *gin.Context) {
	res, err := h.service.GetAll()
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("GetAll roles failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *roleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("GetRoleByID failed: invalid ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	res, err := h.service.GetByID(uint(id))
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("GetRoleByID failed: not found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *roleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Update role failed: invalid ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Update role failed: invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Update role failed: service error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *roleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Delete role failed: invalid ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Delete role failed: service error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
