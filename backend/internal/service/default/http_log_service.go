package service

import (
	"fmt"
	"time"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
)

type HttpLogService interface {
	Create(log *entity.HttpLog) error
	GetAll(query *dto.HttpLogQuery) ([]dto.HttpLogResponse, int64, error)
}

type httpLogService struct {
	repo  repository.HttpLogRepository
	cache cache.CacheService
}

func NewHttpLogService(repo repository.HttpLogRepository, cache cache.CacheService) HttpLogService {
	return &httpLogService{
		repo:  repo,
		cache: cache,
	}
}

func (s *httpLogService) Create(log *entity.HttpLog) error {
	return s.repo.Create(log)
}

func (s *httpLogService) GetAll(query *dto.HttpLogQuery) ([]dto.HttpLogResponse, int64, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf("http_logs:%d:%d:%s:%s:%d",
		query.GetPage(),
		query.GetLimit(),
		query.Method,
		query.Path,
		query.StatusCode,
	)

	type cacheData struct {
		Responses []dto.HttpLogResponse `json:"responses"`
		Total     int64                 `json:"total"`
	}

	var cached cacheData
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return cached.Responses, cached.Total, nil
	}

	logs, total, err := s.repo.FindAll(query)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.HttpLogResponse
	for _, l := range logs {
		responses = append(responses, dto.HttpLogResponse{
			ID:              l.ID,
			RequestID:       l.RequestID,
			Method:          l.Method,
			Path:            l.Path,
			ClientIP:        l.ClientIP,
			UserAgent:       l.UserAgent,
			RequestHeaders:  l.RequestHeaders,
			RequestBody:     l.RequestBody,
			StatusCode:      l.StatusCode,
			ResponseHeaders: l.ResponseHeaders,
			ResponseBody:    l.ResponseBody,
			Latency:         l.Latency,
			UserID:          l.UserID,
			UserEmail:       l.UserEmail,
			MiddlewareTrace: l.MiddlewareTrace,
			CreatedAt:       l.CreatedAt,
		})
	}

	// Save to cache with short TTL (10 seconds)
	_ = s.cache.Set(cacheKey, cacheData{
		Responses: responses,
		Total:     total,
	}, 10*time.Second)

	return responses, total, nil
}
