package service

import (
	"bytes"
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

type MainnnService interface {
	Create(req dto.CreateMainnnRequest) (*dto.MainnnResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.MainnnResponse, error)
	Update(id uint, req dto.UpdateMainnnRequest) (*dto.MainnnResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, error)
}

type mainnnService struct {
	repo  repository.MainnnRepository
	cache cache.CacheService
}

func NewMainnnService(repo repository.MainnnRepository, cache cache.CacheService) MainnnService {
	return &mainnnService{
		repo:  repo,
		cache: cache,
	}
}

func (s *mainnnService) Create(req dto.CreateMainnnRequest) (*dto.MainnnResponse, error) {
	entity := &entity.Mainnn{
		Name:      req.Name,
		Makananan: req.Makananan,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("mainnn:*")

	logger.AuditLogger.Info().
		Uint("mainnn_id", entity.ID).
		Str("action", "mainnn_creation").
		Msg("mainnn created")

	return s.mapToResponse(entity), nil
}

func (s *mainnnService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("mainnn:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.MainnnResponse
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

func (s *mainnnService) GetByID(id uint) (*dto.MainnnResponse, error) {
	cacheKey := fmt.Sprintf("mainnn:%d", id)
	var cached dto.MainnnResponse
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

func (s *mainnnService) Update(id uint, req dto.UpdateMainnnRequest) (*dto.MainnnResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Makananan != "" {
		entity.Makananan = req.Makananan
	}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("mainnn:%d", id))
	s.cache.DeletePattern("mainnn:*")

	logger.AuditLogger.Info().
		Uint("mainnn_id", entity.ID).
		Str("action", "mainnn_update").
		Msg("mainnn updated")

	return s.mapToResponse(entity), nil
}

func (s *mainnnService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("mainnn:%d", id))
	s.cache.DeletePattern("mainnn:*")

	logger.AuditLogger.Info().
		Uint("mainnn_id", id).
		Str("action", "mainnn_deletion").
		Msg("mainnn deleted")

	return s.repo.Delete(id)
}

func (s *mainnnService) Export(format string) ([]byte, error) {
	pagination := &defaultDto.PaginationRequest{
		Page:  1,
		Limit: 1000000,
	}

	entities, _, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	if format == "csv" {
		return s.generateCSV(entities)
	}
	return s.generateExcel(entities)
}

func (s *mainnnService) generateCSV(entities []entity.Mainnn) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "Name", "Makananan", "Created At"}
	writer.Write(header)

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			fmt.Sprintf("%v", e.Name),
			fmt.Sprintf("%v", e.Makananan),
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *mainnnService) generateExcel(entities []entity.Mainnn) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	header := []string{"ID", "Name", "Makananan", "Created At"}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, e := range entities {
		rowNum := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), e.ID)
		cell, _ := excelize.CoordinatesToCellName(0+2, rowNum)
		f.SetCellValue(sheet, cell, e.Name)
		cell, _ = excelize.CoordinatesToCellName(1+2, rowNum)
		f.SetCellValue(sheet, cell, e.Makananan)
		lastCell, _ := excelize.CoordinatesToCellName(len(header), rowNum)
		f.SetCellValue(sheet, lastCell, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *mainnnService) mapToResponse(entity *entity.Mainnn) *dto.MainnnResponse {
	return &dto.MainnnResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		Makananan: entity.Makananan,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
