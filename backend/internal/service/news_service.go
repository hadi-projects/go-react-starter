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

type NewsService interface {
	Create(req dto.CreateNewsRequest) (*dto.NewsResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.NewsResponse, error)
	Update(id uint, req dto.UpdateNewsRequest) (*dto.NewsResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, error)
}

type newsService struct {
	repo  repository.NewsRepository
	cache cache.CacheService
}

func NewNewsService(repo repository.NewsRepository, cache cache.CacheService) NewsService {
	return &newsService{
		repo:  repo,
		cache: cache,
	}
}

func (s *newsService) Create(req dto.CreateNewsRequest) (*dto.NewsResponse, error) {
	entity := &entity.News{
		Name: req.Name,
		Content: req.Content,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("news:*")

	
	logger.AuditLogger.Info().
		Uint("news_id", entity.ID).
		Str("action", "news_creation").
		Msg("news created")
	

	return s.mapToResponse(entity), nil
}

func (s *newsService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("news:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.NewsResponse
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

func (s *newsService) GetByID(id uint) (*dto.NewsResponse, error) {
	cacheKey := fmt.Sprintf("news:%d", id)
	var cached dto.NewsResponse
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

func (s *newsService) Update(id uint, req dto.UpdateNewsRequest) (*dto.NewsResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	entity.Content = req.Content

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("news:%d", id))
	s.cache.DeletePattern("news:*")

	
	logger.AuditLogger.Info().
		Uint("news_id", entity.ID).
		Str("action", "news_update").
		Msg("news updated")
	

	return s.mapToResponse(entity), nil
}

func (s *newsService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("news:%d", id))
	s.cache.DeletePattern("news:*")

	
	logger.AuditLogger.Info().
		Uint("news_id", id).
		Str("action", "news_deletion").
		Msg("news deleted")
	

	return s.repo.Delete(id)
}

func (s *newsService) Export(format string) ([]byte, error) {
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

func (s *newsService) generateCSV(entities []entity.News) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "Name", "Content", "Created At"}
	writer.Write(header)

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			fmt.Sprintf("%v", e.Name),
			fmt.Sprintf("%v", e.Content),
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *newsService) generateExcel(entities []entity.News) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	header := []string{"ID", "Name", "Content", "Created At"}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, e := range entities {
		rowNum := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), e.ID)
		
		var cell string
		cell, _ = excelize.CoordinatesToCellName(0+2, rowNum)
		f.SetCellValue(sheet, cell, e.Name)
		
		cell, _ = excelize.CoordinatesToCellName(1+2, rowNum)
		f.SetCellValue(sheet, cell, e.Content)
		
		lastCell, _ := excelize.CoordinatesToCellName(len(header), rowNum)
		f.SetCellValue(sheet, lastCell, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *newsService) mapToResponse(entity *entity.News) *dto.NewsResponse {
	return &dto.NewsResponse{
		ID:        entity.ID,
		Name: entity.Name,
		Content: entity.Content,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
