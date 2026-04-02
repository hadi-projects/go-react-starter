package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type RoleHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Export(c *gin.Context)
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create role failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Role created successfully", res)
}

func (h *roleHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")
	category := c.DefaultQuery("category", "")

	pagination := &dto.PaginationRequest{
		Page:     page,
		Limit:    limit,
		Search:   search,
		Category: category,
	}

	res, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("GetAll roles failed")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &response.PaginationMeta{
		CurrentPage: res.Meta.CurrentPage,
		TotalPages:  res.Meta.TotalPages,
		TotalData:   res.Meta.TotalData,
		Limit:       res.Meta.Limit,
	}

	response.SuccessWithPagination(c, http.StatusOK, "Roles retrieved successfully", res.Data, meta)
}

func (h *roleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("GetRoleByID failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("GetRoleByID failed: not found")
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Role retrieved successfully", res)
}

func (h *roleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Update role failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Update role failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Update role failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Role updated successfully", res)
}

func (h *roleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Delete role failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Delete role failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Role deleted successfully", nil)
}

func (h *roleHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "excel")
	if format != "csv" && format != "excel" {
		response.Error(c, http.StatusBadRequest, "Invalid format. Supported: csv, excel")
		return
	}

	data, filename, err := h.service.Export(c.Request.Context(), format)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Export roles failed")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if format == "csv" {
		contentType = "text/csv"
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
}
