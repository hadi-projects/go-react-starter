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

type AuditLogService interface {
	Create(log *entity.AuditLog) error
	GetAll(query *dto.AuditLogQuery) ([]dto.AuditLogResponse, int64, error)
	Export(query *dto.AuditLogQuery, format string) ([]byte, string, error)
}

type auditLogService struct {
	repo  repository.AuditLogRepository
	cache cache.CacheService
}

func NewAuditLogService(repo repository.AuditLogRepository, cache cache.CacheService) AuditLogService {
	return &auditLogService{
		repo:  repo,
		cache: cache,
	}
}

func (s *auditLogService) Create(log *entity.AuditLog) error {
	return s.repo.Create(&logger.AuditLog{
		RequestID: log.RequestID,
		UserID:    log.UserID,
		UserEmail: log.UserEmail,
		Action:    log.Action,
		Module:    log.Module,
		TargetID:  log.TargetID,
		Metadata:  log.Metadata,
	})
}

func (s *auditLogService) GetAll(query *dto.AuditLogQuery) ([]dto.AuditLogResponse, int64, error) {
	cacheKey := fmt.Sprintf("audit_logs:%d:%d:%s:%s:%s:%s",
		query.GetPage(),
		query.GetLimit(),
		query.Module,
		query.Action,
		query.UserEmail,
		query.RequestID,
	)

	type cacheData struct {
		Responses []dto.AuditLogResponse `json:"responses"`
		Total     int64                  `json:"total"`
	}

	var cached cacheData
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return cached.Responses, cached.Total, nil
	}

	logs, total, err := s.repo.FindAll(query)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.AuditLogResponse
	for _, l := range logs {
		responses = append(responses, dto.AuditLogResponse{
			ID:        l.ID,
			RequestID: l.RequestID,
			UserID:    l.UserID,
			UserEmail: l.UserEmail,
			Action:    l.Action,
			Module:    l.Module,
			TargetID:  l.TargetID,
			Metadata:  l.Metadata,
			CreatedAt: l.CreatedAt,
		})
	}

	_ = s.cache.Set(cacheKey, cacheData{
		Responses: responses,
		Total:     total,
	}, 10*time.Second)

	return responses, total, nil
}

func (s *auditLogService) Export(query *dto.AuditLogQuery, format string) ([]byte, string, error) {
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

func (s *auditLogService) generateCSV(logs []entity.AuditLog) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Time", "Request ID", "User ID", "Email", "Action", "Module", "Target ID", "Metadata"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, l := range logs {
		row := []string{
			fmt.Sprintf("%d", l.ID),
			l.CreatedAt.Format("2006-01-02 15:04:05"),
			l.RequestID,
			fmt.Sprintf("%d", l.UserID),
			l.UserEmail,
			l.Action,
			l.Module,
			l.TargetID,
			l.Metadata,
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "audit_logs.csv", nil
}

func (s *auditLogService) generateExcel(logs []entity.AuditLog) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Audit Logs"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Time", "Request ID", "User ID", "Email", "Action", "Module", "Target ID", "Metadata"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, l := range logs {
		row := []interface{}{
			l.ID,
			l.CreatedAt.Format("2006-01-02 15:04:05"),
			l.RequestID,
			l.UserID,
			l.UserEmail,
			l.Action,
			l.Module,
			l.TargetID,
			l.Metadata,
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

	return buf.Bytes(), "audit_logs.xlsx", nil
}
