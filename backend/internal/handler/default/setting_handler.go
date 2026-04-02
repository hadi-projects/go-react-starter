package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type SettingHandler interface {
	GetByCategory(c *gin.Context)
	BulkUpdate(c *gin.Context)
}

type settingHandler struct {
	service service.SettingService
}

func NewSettingHandler(service service.SettingService) SettingHandler {
	return &settingHandler{service: service}
}

func (h *settingHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	res, err := h.service.GetSettings(c.Request.Context(), category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Settings retrieved", res)
}

func (h *settingHandler) BulkUpdate(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.BulkUpdate(c.Request.Context(), req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Settings updated successfully", nil)
}
