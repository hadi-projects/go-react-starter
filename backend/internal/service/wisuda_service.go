package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	"github.com/xuri/excelize/v2"
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type WisudaService interface {
	Create(req dto.CreateWisudaRequest) (*dto.WisudaResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.WisudaResponse, error)
	Update(id uint, req dto.UpdateWisudaRequest) (*dto.WisudaResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, error)
}

type wisudaService struct {
	repo  repository.WisudaRepository
	cache cache.CacheService
}

func NewWisudaService(repo repository.WisudaRepository, cache cache.CacheService) WisudaService {
	return &wisudaService{
		repo:  repo,
		cache: cache,
	}
}

func (s *wisudaService) Create(req dto.CreateWisudaRequest) (*dto.WisudaResponse, error) {
	entity := &entity.Wisuda{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("wisuda:*")

	
	logger.AuditLogger.Info().
		Uint("wisuda_id", entity.ID).
		Str("action", "wisuda_creation").
		Msg("wisuda created")
	

	return s.mapToResponse(entity), nil
}

func (s *wisudaService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("wisuda:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.WisudaResponse
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

func (s *wisudaService) GetByID(id uint) (*dto.WisudaResponse, error) {
	cacheKey := fmt.Sprintf("wisuda:%d", id)
	var cached dto.WisudaResponse
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

func (s *wisudaService) Update(id uint, req dto.UpdateWisudaRequest) (*dto.WisudaResponse, error) {
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

	s.cache.Delete(fmt.Sprintf("wisuda:%d", id))
	s.cache.DeletePattern("wisuda:*")

	
	logger.AuditLogger.Info().
		Uint("wisuda_id", entity.ID).
		Str("action", "wisuda_update").
		Msg("wisuda updated")
	

	return s.mapToResponse(entity), nil
}

func (s *wisudaService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("wisuda:%d", id))
	s.cache.DeletePattern("wisuda:*")

	
	logger.AuditLogger.Info().
		Uint("wisuda_id", id).
		Str("action", "wisuda_deletion").
		Msg("wisuda deleted")
	

	return s.repo.Delete(id)
}

func (s *wisudaService) Export(format string) ([]byte, error) {
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

func (s *wisudaService) generateCSV(entities []entity.Wisuda) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "Name", "Created At"}
	writer.Write(header)

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			fmt.Sprintf("%v", e.Name),
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *wisudaService) generateExcel(entities []entity.Wisuda) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	header := []string{"ID", "Name", "Created At"}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, e := range entities {
		rowNum := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), e.ID)
		cell, _ := excelize.CoordinatesToCellName(0+2, rowNum)
		f.SetCellValue(sheet, cell, e.Name)
		lastCell, _ := excelize.CoordinatesToCellName(len(header), rowNum)
		f.SetCellValue(sheet, lastCell, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *wisudaService) mapToResponse(entity *entity.Wisuda) *dto.WisudaResponse {
	return &dto.WisudaResponse{
		ID:        entity.ID,
		Name: entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
