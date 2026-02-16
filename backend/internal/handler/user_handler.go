package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/service"
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Register(req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Register failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", res)
}

func (h *userHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create user failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateUser(req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Create user failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", res)
}

func (h *userHandler) Me(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		logger.SystemLogger.Error().Msg("Me failed: user_id not found in context")
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := val.(uint)
	if !ok {
		logger.SystemLogger.Error().Msg("Me failed: invalid user_id type")
		response.Error(c, http.StatusInternalServerError, "Invalid user ID type")
		return
	}

	res, err := h.service.GetMe(userID)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("user_id", userID).Msg("Me failed: user not found")
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved successfully", res)
}

func (h *userHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	pagination := &dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	res, err := h.service.GetAll(pagination)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("GetAll users failed")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &response.PaginationMeta{
		CurrentPage: res.Meta.CurrentPage,
		TotalPages:  res.Meta.TotalPages,
		TotalItems:  res.Meta.TotalItems,
		Limit:       res.Meta.Limit,
	}

	response.SuccessWithPagination(c, http.StatusOK, "Users retrieved successfully", res.Data, meta)
}

func (h *userHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Update user failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Update user failed: invalid request body")
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Update user failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "User updated successfully", res)
}

func (h *userHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Str("id", idStr).Msg("Delete user failed: invalid ID")
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		logger.SystemLogger.Error().Err(err).Uint("id", uint(id)).Msg("Delete user failed: service error")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
