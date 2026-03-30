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

type PermissionHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Export(c *gin.Context)
}

type permissionHandler struct {
	service service.PermissionService
}

func NewPermissionHandler(service service.PermissionService) PermissionHandler {
	return &permissionHandler{service: service}
}

func (h *permissionHandler) Create(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create permission failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create permission failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Permission created successfully", res)
}

func (h *permissionHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	pagination := &dto.PaginationRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	res, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("GetAll permissions failed")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &response.PaginationMeta{
		CurrentPage: res.Meta.CurrentPage,
		TotalPages:  res.Meta.TotalPages,
		TotalData:   res.Meta.TotalData,
		Limit:       res.Meta.Limit,
	}

	response.SuccessWithPagination(c, http.StatusOK, "Permissions retrieved successfully", res.Data, meta)
}

func (h *permissionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Update permission failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Update permission failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Update permission failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permission updated successfully", res)
}

func (h *permissionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Delete permission failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Delete permission failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permission deleted successfully", nil)
}

func (h *permissionHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "excel")
	if format != "csv" && format != "excel" {
		response.Error(c, http.StatusBadRequest, "Invalid format. Supported: csv, excel")
		return
	}

	data, filename, err := h.service.Export(c.Request.Context(), format)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Export permissions failed")
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
