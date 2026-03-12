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
	"github.com/xuri/excelize/v2"
)

type HttpLogService interface {
	Create(log *entity.HttpLog) error
	GetAll(query *dto.HttpLogQuery) ([]dto.HttpLogResponse, int64, error)
	Export(query *dto.HttpLogQuery, format string) ([]byte, string, error)
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

func (s *httpLogService) Export(query *dto.HttpLogQuery, format string) ([]byte, string, error) {
	// For export, we want to fetch all matching logs without small pagination limits
	// But let's reuse FindAll with a large limit
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

func (s *httpLogService) generateCSV(logs []entity.HttpLog) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Time", "Request ID", "Method", "Path", "Status", "Latency", "User Email"}
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
			l.UserEmail,
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "http_logs.csv", nil
}

func (s *httpLogService) generateExcel(logs []entity.HttpLog) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "HTTP Logs"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Time", "Request ID", "Method", "Path", "Status", "Latency", "User Email"}
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
			l.UserEmail,
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

	return buf.Bytes(), "http_logs.xlsx", nil
}
