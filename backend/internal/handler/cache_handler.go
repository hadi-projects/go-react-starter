package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type CacheHandler interface {
	ClearAll(c *gin.Context)
}

type cacheHandler struct {
	cache cache.CacheService
}

func NewCacheHandler(cache cache.CacheService) CacheHandler {
	return &cacheHandler{cache: cache}
}

func (h *cacheHandler) ClearAll(c *gin.Context) {
	if err := h.cache.FlushAll(); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Failed to clear cache")
		response.Error(c, http.StatusInternalServerError, "Failed to clear cache")
		return
	}

	logger.SystemLogger.Info().Msg("Cache cleared successfully")
	response.Success(c, http.StatusOK, "Cache cleared successfully", nil)
}
