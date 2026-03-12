package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type TestsajaService interface {
	Create(ctx context.Context, req dto.CreateTestsajaRequest) (*dto.TestsajaResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.TestsajaResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateTestsajaRequest) (*dto.TestsajaResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, format string) ([]byte, string, error)
}

type testsajaService struct {
	repo  repository.TestsajaRepository
	cache cache.CacheService
}

func NewTestsajaService(repo repository.TestsajaRepository, cache cache.CacheService) TestsajaService {
	return &testsajaService{
		repo:  repo,
		cache: cache,
	}
}

func (s *testsajaService) Create(ctx context.Context, req dto.CreateTestsajaRequest) (*dto.TestsajaResponse, error) {
	entity := &entity.Testsaja{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("testsaja:*")

	// logger.AuditLogger.Info().
	// 	Uint("testsaja_id", entity.ID).
	// 	Str("action", "testsaja_creation").
	// 	Msg("testsaja created")
	logger.LogAudit(ctx, "CREATE", "TESTSAJA", fmt.Sprintf("%d", entity.ID), fmt.Sprintf("name: %s", entity.Name))

	return s.mapToResponse(entity), nil
}

func (s *testsajaService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("testsaja:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.TestsajaResponse
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

func (s *testsajaService) GetByID(id uint) (*dto.TestsajaResponse, error) {
	cacheKey := fmt.Sprintf("testsaja:%d", id)
	var cached dto.TestsajaResponse
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

func (s *testsajaService) Update(ctx context.Context, id uint, req dto.UpdateTestsajaRequest) (*dto.TestsajaResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("testsaja:%d", id))
	s.cache.DeletePattern("testsaja:*")

	// logger.AuditLogger.Info().
	// 	Uint("testsaja_id", entity.ID).
	// 	Str("action", "testsaja_update").
	// 	Msg("testsaja updated")
	logger.LogAudit(ctx, "UPDATE", "TESTSAJA", fmt.Sprintf("%d", id), fmt.Sprintf("name: %s", entity.Name))

	return s.mapToResponse(entity), nil
}

func (s *testsajaService) Delete(ctx context.Context, id uint) error {
	s.cache.Delete(fmt.Sprintf("testsaja:%d", id))
	s.cache.DeletePattern("testsaja:*")

	// logger.AuditLogger.Info().
	// 	Uint("testsaja_id", id).
	// 	Str("action", "testsaja_deletion").
	// 	Msg("testsaja deleted")
	logger.LogAudit(ctx, "DELETE", "TESTSAJA", fmt.Sprintf("%d", id), "")

	return s.repo.Delete(id)
}

func (s *testsajaService) Export(ctx context.Context, format string) ([]byte, string, error) {
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

func (s *testsajaService) generateCSV(entities []entity.Testsaja) ([]byte, string, error) {
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
	return buf.Bytes(), "testsaja.csv", nil
}

func (s *testsajaService) generateExcel(entities []entity.Testsaja) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Testsaja"
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

	return buf.Bytes(), "testsaja.xlsx", nil
}

func (s *testsajaService) mapToResponse(entity *entity.Testsaja) *dto.TestsajaResponse {
	return &dto.TestsajaResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
