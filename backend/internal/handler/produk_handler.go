package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/service"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type ProdukHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Export(c *gin.Context)
}

type produkHandler struct {
	service service.ProdukService
}

func NewProdukHandler(service service.ProdukService) ProdukHandler {
	return &produkHandler{service: service}
}

func (h *produkHandler) Create(c *gin.Context) {
	var req dto.CreateProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Produk created successfully", res)
}

func (h *produkHandler) GetAll(c *gin.Context) {
	var pagination defaultDto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.GetAll(&pagination)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Produks retrieved successfully", res)
}

func (h *produkHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Produk not found")
		return
	}

	response.Success(c, http.StatusOK, "Produk retrieved successfully", res)
}

func (h *produkHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Produk updated successfully", res)
}

func (h *produkHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Produk deleted successfully", nil)
}

func (h *produkHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "excel")
	if format != "csv" && format != "excel" {
		response.Error(c, http.StatusBadRequest, "Invalid format. Supported: csv, excel")
		return
	}

	data, filename, err := h.service.Export(c.Request.Context(), format)
	if err != nil {
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
