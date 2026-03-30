package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type HttpLogHandler interface {
	GetAll(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type httpLogHandler struct {
	httpLogService service.HttpLogService
}

func NewHttpLogHandler(httpLogService service.HttpLogService) HttpLogHandler {
	return &httpLogHandler{httpLogService: httpLogService}
}

func (h *httpLogHandler) GetAll(ctx *gin.Context) {
	var query dto.HttpLogQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	logs, total, err := h.httpLogService.GetAll(ctx.Request.Context(), &query)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	paginationMeta := &response.PaginationMeta{
		CurrentPage: query.GetPage(),
		TotalPages:  int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
		TotalData:   total,
		Limit:       query.GetLimit(),
	}

	response.SuccessWithPagination(ctx, http.StatusOK, "HTTP logs retrieved successfully", logs, paginationMeta)
}

func (h *httpLogHandler) Export(ctx *gin.Context) {
	var query dto.HttpLogQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	format := ctx.DefaultQuery("format", "excel")
	data, filename, err := h.httpLogService.Export(ctx.Request.Context(), &query, format)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if format == "csv" {
		contentType = "text/csv"
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", contentType)
	ctx.Data(http.StatusOK, contentType, data)
}
