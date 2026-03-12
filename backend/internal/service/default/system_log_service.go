package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type SystemLogService interface {
	Create(log *entity.SystemLog) error
	GetAll(query *dto.SystemLogQuery) ([]dto.SystemLogResponse, int64, error)
	Export(query *dto.SystemLogQuery, format string) ([]byte, string, error)
}

type systemLogService struct {
	repo  repository.SystemLogRepository
	cache cache.CacheService
}

func NewSystemLogService(repo repository.SystemLogRepository, cache cache.CacheService) SystemLogService {
	return &systemLogService{
		repo:  repo,
		cache: cache,
	}
}

func (s *systemLogService) Create(log *entity.SystemLog) error {
	return s.repo.Create(&logger.SystemLog{
		RequestID:    log.RequestID,
		Method:       log.Method,
		Path:         log.Path,
		StatusCode:   log.StatusCode,
		Latency:      log.Latency,
		RequestBody:  log.RequestBody,
		ResponseBody: log.ResponseBody,
	})
}

func (s *systemLogService) GetAll(query *dto.SystemLogQuery) ([]dto.SystemLogResponse, int64, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf("system_logs:%d:%d:%s:%s:%d:%s",
		query.GetPage(),
		query.GetLimit(),
		query.Method,
		query.Path,
		query.StatusCode,
		query.RequestID,
	)

	type cacheData struct {
		Responses []dto.SystemLogResponse `json:"responses"`
		Total     int64                   `json:"total"`
	}

	var cached cacheData
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return cached.Responses, cached.Total, nil
	}

	logs, total, err := s.repo.FindAll(query)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.SystemLogResponse
	for _, l := range logs {
		responses = append(responses, dto.SystemLogResponse{
			ID:           l.ID,
			RequestID:    l.RequestID,
			Method:       l.Method,
			Path:         l.Path,
			StatusCode:   l.StatusCode,
			Latency:      l.Latency,
			RequestBody:  l.RequestBody,
			ResponseBody: l.ResponseBody,
			CreatedAt:    l.CreatedAt,
		})
	}

	// Save to cache with short TTL (10 seconds)
	_ = s.cache.Set(cacheKey, cacheData{
		Responses: responses,
		Total:     total,
	}, 10*time.Second)

	return responses, total, nil
}

func (s *systemLogService) Export(query *dto.SystemLogQuery, format string) ([]byte, string, error) {
	exportQuery := *query
	exportQuery.Limit = 1000000

	logs, _, err := s.repo.FindAll(&exportQuery)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(logs)
	}
	return s.generateExcel(logs)
}

func (s *systemLogService) generateCSV(logs []entity.SystemLog) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Time", "Request ID", "Method", "Path", "Status", "Latency"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, l := range logs {
		row := []string{
			fmt.Sprintf("%d", l.ID),
			l.CreatedAt.Format("2006-01-02 15:04:05"),
			l.RequestID,
			l.Method,
			l.Path,
			fmt.Sprintf("%d", l.StatusCode),
			fmt.Sprintf("%dms", l.Latency),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "system_logs.csv", nil
}

func (s *systemLogService) generateExcel(logs []entity.SystemLog) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "System Logs"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Time", "Request ID", "Method", "Path", "Status", "Latency"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, l := range logs {
		row := []interface{}{
			l.ID,
			l.CreatedAt.Format("2006-01-02 15:04:05"),
			l.RequestID,
			l.Method,
			l.Path,
			l.StatusCode,
			l.Latency,
		}
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), "system_logs.xlsx", nil
}
