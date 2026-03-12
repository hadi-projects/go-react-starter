package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type TestduaService interface {
	Create(ctx context.Context, req dto.CreateTestduaRequest) (*dto.TestduaResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.TestduaResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateTestduaRequest) (*dto.TestduaResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, format string) ([]byte, string, error)
}

type testduaService struct {
	repo  repository.TestduaRepository
	cache cache.CacheService
}

func NewTestduaService(repo repository.TestduaRepository, cache cache.CacheService) TestduaService {
	return &testduaService{
		repo:  repo,
		cache: cache,
	}
}

func (s *testduaService) Create(ctx context.Context, req dto.CreateTestduaRequest) (*dto.TestduaResponse, error) {
	entity := &entity.Testdua{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("testdua:*")

	
	// logger.AuditLogger.Info().
	// 	Uint("testdua_id", entity.ID).
	// 	Str("action", "testdua_creation").
	// 	Msg("testdua created")
	logger.LogAudit(ctx, "CREATE", "TESTDUA", fmt.Sprintf("%d", entity.ID), fmt.Sprintf("name: %s", entity.Name))
	

	return s.mapToResponse(entity), nil
}

func (s *testduaService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("testdua:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.TestduaResponse
	for _, e := range entities {
		responses = append(responses, *s.mapToResponse(&e))
	}

	response := &defaultDto.PaginationResponse{
		Data: responses,
		Meta: defaultDto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  int(math.Ceil(float64(total) / float64(pagination.GetLimit()))),
			TotalData:   total,
			Limit:       pagination.GetLimit(),
		},
	}

	s.cache.Set(cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *testduaService) GetByID(id uint) (*dto.TestduaResponse, error) {
	cacheKey := fmt.Sprintf("testdua:%d", id)
	var cached dto.TestduaResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.mapToResponse(entity)
	s.cache.Set(cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *testduaService) Update(ctx context.Context, id uint, req dto.UpdateTestduaRequest) (*dto.TestduaResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	entity.Name = req.Name

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("testdua:%d", id))
	s.cache.DeletePattern("testdua:*")

	
	// logger.AuditLogger.Info().
	// 	Uint("testdua_id", entity.ID).
	// 	Str("action", "testdua_update").
	// 	Msg("testdua updated")
	logger.LogAudit(ctx, "UPDATE", "TESTDUA", fmt.Sprintf("%d", id), fmt.Sprintf("name: %s", entity.Name))
	

	return s.mapToResponse(entity), nil
}

func (s *testduaService) Delete(ctx context.Context, id uint) error {
	s.cache.Delete(fmt.Sprintf("testdua:%d", id))
	s.cache.DeletePattern("testdua:*")

	
	// logger.AuditLogger.Info().
	// 	Uint("testdua_id", id).
	// 	Str("action", "testdua_deletion").
	// 	Msg("testdua deleted")
	logger.LogAudit(ctx, "DELETE", "TESTDUA", fmt.Sprintf("%d", id), "")
	

	return s.repo.Delete(id)
}

func (s *testduaService) Export(ctx context.Context, format string) ([]byte, string, error) {
	pagination := &defaultDto.PaginationRequest{
		Page:  1,
		Limit: 1000000,
	}

	entities, _, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(entities)
	}
	return s.generateExcel(entities)
}

func (s *testduaService) generateCSV(entities []entity.Testdua) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Name", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			e.Name,
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "testdua.csv", nil
}

func (s *testduaService) generateExcel(entities []entity.Testdua) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Testdua"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Name", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, e := range entities {
		row := []interface{}{
			e.ID,
			e.Name,
			e.CreatedAt.Format("2006-01-02 15:04:05"),
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

	return buf.Bytes(), "testdua.xlsx", nil
}

func (s *testduaService) mapToResponse(entity *entity.Testdua) *dto.TestduaResponse {
	return &dto.TestduaResponse{
		ID:        entity.ID,
		Name: entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
