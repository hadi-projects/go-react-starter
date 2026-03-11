package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type HttpLogHandler interface {
	GetAll(ctx *gin.Context)
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

	logs, total, err := h.httpLogService.GetAll(&query)
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
