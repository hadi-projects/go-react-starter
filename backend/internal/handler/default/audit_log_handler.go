package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
)

type AuditLogHandler interface {
	GetAll(c *gin.Context)
	Export(c *gin.Context)
}

type auditLogHandler struct {
	service service.AuditLogService
}

func NewAuditLogHandler(service service.AuditLogService) AuditLogHandler {
	return &auditLogHandler{service: service}
}

func (h *auditLogHandler) GetAll(c *gin.Context) {
	var query dto.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, total, err := h.service.GetAll(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  query.GetPage(),
		"limit": query.GetLimit(),
	})
}

func (h *auditLogHandler) Export(c *gin.Context) {
	var query dto.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	format := c.DefaultQuery("format", "excel")
	data, filename, err := h.service.Export(&query, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
