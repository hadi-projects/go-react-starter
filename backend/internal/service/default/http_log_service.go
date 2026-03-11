package service

import (
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
)

type HttpLogService interface {
	Create(log *entity.HttpLog) error
	GetAll(query *dto.HttpLogQuery) ([]dto.HttpLogResponse, int64, error)
}

type httpLogService struct {
	repo repository.HttpLogRepository
}

func NewHttpLogService(repo repository.HttpLogRepository) HttpLogService {
	return &httpLogService{repo: repo}
}

func (s *httpLogService) Create(log *entity.HttpLog) error {
	return s.repo.Create(log)
}

func (s *httpLogService) GetAll(query *dto.HttpLogQuery) ([]dto.HttpLogResponse, int64, error) {
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
			CreatedAt:       l.CreatedAt,
		})
	}

	return responses, total, nil
}
