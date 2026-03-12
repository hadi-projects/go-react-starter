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

type UserHandler interface {
	Register(c *gin.Context)
	Create(c *gin.Context)
	Me(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Export(c *gin.Context)
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
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Register failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Register failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", res)
}

func (h *userHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Create user failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Create user failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", res)
}

func (h *userHandler) Me(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		logger.WithCtx(c, logger.SystemLogger).Error().Msg("Me failed: user_id not found in context")
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := val.(uint)
	if !ok {
		logger.WithCtx(c, logger.SystemLogger).Error().Msg("Me failed: invalid user_id type")
		response.Error(c, http.StatusInternalServerError, "Invalid user ID type")
		return
	}

	res, err := h.service.GetMe(userID)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Uint("user_id", userID).Msg("Me failed: user not found")
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved successfully", res)
}

func (h *userHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	pagination := &dto.PaginationRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	res, err := h.service.GetAll(pagination)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("GetAll users failed")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &response.PaginationMeta{
		CurrentPage: res.Meta.CurrentPage,
		TotalPages:  res.Meta.TotalPages,
		TotalData:   res.Meta.TotalData,
		Limit:       res.Meta.Limit,
	}

	response.SuccessWithPagination(c, http.StatusOK, "Users retrieved successfully", res.Data, meta)
}

func (h *userHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Str("id", idStr).Msg("Update user failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Update user failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Uint("id", uint(id)).Msg("Update user failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "User updated successfully", res)
}

func (h *userHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Str("id", idStr).Msg("Delete user failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Uint("id", uint(id)).Msg("Delete user failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *userHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "excel")
	if format != "csv" && format != "excel" {
		response.Error(c, http.StatusBadRequest, "Invalid format. Supported: csv, excel")
		return
	}

	data, filename, err := h.service.Export(c.Request.Context(), format)
	if err != nil {
		logger.WithCtx(c, logger.SystemLogger).Error().Err(err).Msg("Export users failed")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if format == "csv" {
		contentType = "text/csv"
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Data(http.StatusOK, contentType, data)
}
