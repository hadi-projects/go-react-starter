package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type ApiKeyHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Delete(c *gin.Context)
}

type apiKeyHandler struct {
	service service.ApiKeyService
}

func NewApiKeyHandler(service service.ApiKeyService) ApiKeyHandler {
	return &apiKeyHandler{service: service}
}

func (h *apiKeyHandler) Create(c *gin.Context) {
	var req dto.ApiKeyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	userID := userIDValue.(uint)

	res, err := h.service.Generate(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "API Key generated successfully", res)
}

func (h *apiKeyHandler) GetAll(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		// Fallback to defaults
		pagination.Page = 1
		pagination.Limit = 10
	}

	keys, total, err := h.service.GetAll(c.Request.Context(), &pagination)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := dto.PaginationMeta{
		CurrentPage: pagination.GetPage(),
		TotalPages:  int((total + int64(pagination.GetLimit()) - 1) / int64(pagination.GetLimit())),
		TotalData:   total,
		Limit:       pagination.GetLimit(),
	}

	response.Success(c, http.StatusOK, "API Keys retrieved", gin.H{
		"data": keys,
		"meta": meta,
	})
}

func (h *apiKeyHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "API Key revoked successfully", nil)
}
