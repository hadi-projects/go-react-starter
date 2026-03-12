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

type ArsipService interface {
	Create(req dto.CreateArsipRequest) (*dto.ArsipResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.ArsipResponse, error)
	Update(id uint, req dto.UpdateArsipRequest) (*dto.ArsipResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, error)
}

type arsipService struct {
	repo  repository.ArsipRepository
	cache cache.CacheService
}

func NewArsipService(repo repository.ArsipRepository, cache cache.CacheService) ArsipService {
	return &arsipService{
		repo:  repo,
		cache: cache,
	}
}

func (s *arsipService) Create(req dto.CreateArsipRequest) (*dto.ArsipResponse, error) {
	entity := &entity.Arsip{
		Name:    req.Name,
		Tanggal: req.Tanggal,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("arsip:*")

	logger.AuditLogger.Info().
		Uint("arsip_id", entity.ID).
		Str("action", "arsip_creation").
		Msg("arsip created")

	return s.mapToResponse(entity), nil
}

func (s *arsipService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("arsip:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.ArsipResponse
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

func (s *arsipService) GetByID(id uint) (*dto.ArsipResponse, error) {
	cacheKey := fmt.Sprintf("arsip:%d", id)
	var cached dto.ArsipResponse
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

func (s *arsipService) Update(id uint, req dto.UpdateArsipRequest) (*dto.ArsipResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Tanggal != "" {
		entity.Tanggal = req.Tanggal
	}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("arsip:%d", id))
	s.cache.DeletePattern("arsip:*")

	logger.AuditLogger.Info().
		Uint("arsip_id", entity.ID).
		Str("action", "arsip_update").
		Msg("arsip updated")

	return s.mapToResponse(entity), nil
}

func (s *arsipService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("arsip:%d", id))
	s.cache.DeletePattern("arsip:*")

	logger.AuditLogger.Info().
		Uint("arsip_id", id).
		Str("action", "arsip_deletion").
		Msg("arsip deleted")

	return s.repo.Delete(id)
}

func (s *arsipService) Export(format string) ([]byte, error) {
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

func (s *arsipService) generateCSV(entities []entity.Arsip) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "Name", "Tanggal", "Created At"}
	writer.Write(header)

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			fmt.Sprintf("%v", e.Name),
			fmt.Sprintf("%v", e.Tanggal),
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *arsipService) generateExcel(entities []entity.Arsip) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	header := []string{"ID", "Name", "Tanggal", "Created At"}
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
		f.SetCellValue(sheet, cell, e.Tanggal)
		lastCell, _ := excelize.CoordinatesToCellName(len(header), rowNum)
		f.SetCellValue(sheet, lastCell, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *arsipService) mapToResponse(entity *entity.Arsip) *dto.ArsipResponse {
	return &dto.ArsipResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		Tanggal:   entity.Tanggal,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
