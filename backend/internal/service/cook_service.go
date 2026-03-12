package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	defaultdto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type CookService interface {
	Create(req dto.CreateCookRequest) (*dto.CookResponse, error)
	GetAll(pagination *defaultdto.PaginationRequest) (*defaultdto.PaginationResponse, error)
	GetByID(id uint) (*dto.CookResponse, error)
	Update(id uint, req dto.UpdateCookRequest) (*dto.CookResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, string, error)
}

type cookService struct {
	repo  repository.CookRepository
	cache cache.CacheService
}

func NewCookService(repo repository.CookRepository, cache cache.CacheService) CookService {
	return &cookService{
		repo:  repo,
		cache: cache,
	}
}

func (s *cookService) Create(req dto.CreateCookRequest) (*dto.CookResponse, error) {
	entity := &entity.Cook{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("cook:*")

	
	logger.AuditLogger.Info().
		Uint("cook_id", entity.ID).
		Str("action", "cook_creation").
		Msg("cook created")
	

	return s.mapToResponse(entity), nil
}

func (s *cookService) GetAll(pagination *defaultdto.PaginationRequest) (*defaultdto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("cook:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultdto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.CookResponse
	for _, e := range entities {
		responses = append(responses, *s.mapToResponse(&e))
	}

	response := &defaultdto.PaginationResponse{
		Data: responses,
		Meta: defaultdto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  int(math.Ceil(float64(total) / float64(pagination.GetLimit()))),
			TotalData:   total,
			Limit:       pagination.GetLimit(),
		},
	}

	s.cache.Set(cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *cookService) GetByID(id uint) (*dto.CookResponse, error) {
	cacheKey := fmt.Sprintf("cook:%d", id)
	var cached dto.CookResponse
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

func (s *cookService) Update(id uint, req dto.UpdateCookRequest) (*dto.CookResponse, error) {
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

	s.cache.Delete(fmt.Sprintf("cook:%d", id))
	s.cache.DeletePattern("cook:*")

	
	logger.AuditLogger.Info().
		Uint("cook_id", entity.ID).
		Str("action", "cook_update").
		Msg("cook updated")
	

	return s.mapToResponse(entity), nil
}

func (s *cookService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("cook:%d", id))
	s.cache.DeletePattern("cook:*")

	
	logger.AuditLogger.Info().
		Uint("cook_id", id).
		Str("action", "cook_deletion").
		Msg("cook deleted")
	

	return s.repo.Delete(id)
}

func (s *cookService) Export(format string) ([]byte, string, error) {
	pagination := &defaultdto.PaginationRequest{
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

func (s *cookService) generateCSV(entities []entity.Cook) ([]byte, string, error) {
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
	return buf.Bytes(), "cook.csv", nil
}

func (s *cookService) generateExcel(entities []entity.Cook) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Cook"
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

	return buf.Bytes(), "cook.xlsx", nil
}

func (s *cookService) mapToResponse(entity *entity.Cook) *dto.CookResponse {
	return &dto.CookResponse{
		ID:        entity.ID,
		Name: entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
