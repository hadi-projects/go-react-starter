package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/kafka"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type HealthHandler interface {
	GetStatus(c *gin.Context)
}

type healthHandler struct {
	cache    cache.CacheService
	producer kafka.Producer
}

func NewHealthHandler(cache cache.CacheService, producer kafka.Producer) HealthHandler {
	return &healthHandler{
		cache:    cache,
		producer: producer,
	}
}

func (h *healthHandler) GetStatus(c *gin.Context) {
	redisStatus := h.cache.Status()
	
	kafkaStatus := "disconnected"
	if h.producer != nil {
		kafkaStatus = h.producer.Status()
	}

	response.Success(c, http.StatusOK, "System status", gin.H{
		"redis": redisStatus,
		"kafka": kafkaStatus,
	})
}
