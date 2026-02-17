package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type LogHandler interface {
	GetLogs(ctx *gin.Context)
}

type logHandler struct {
	logService service.LogService
}

func NewLogHandler(logService service.LogService) LogHandler {
	return &logHandler{logService: logService}
}

func (h *logHandler) GetLogs(ctx *gin.Context) {
	var query dto.LogQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Check permissions
	permissions, exists := ctx.Get("permissions")
	if !exists {
		response.Error(ctx, http.StatusForbidden, "unauthorized")
		return
	}

	perms := permissions.([]string)
	isAdmin := false
	canSeeOwn := false
	for _, p := range perms {
		if p == "get-all-logs" {
			isAdmin = true
			break
		}
		if p == "get-own-logs" {
			canSeeOwn = true
		}
	}

	if !isAdmin && !canSeeOwn {
		response.Error(ctx, http.StatusForbidden, "you don't have permission to view logs")
		return
	}

	// If not admin, force filter by user_id
	if !isAdmin {
		userID, exists := ctx.Get("user_id")
		if !exists {
			response.Error(ctx, http.StatusForbidden, "user id not found in context")
			return
		}
		query.UserID = userID.(uint)
	}

	logs, err := h.logService.GetLogs(query)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Logs retrieved successfully", logs)
}
