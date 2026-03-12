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

type BlogService interface {
	Create(req dto.CreateBlogRequest) (*dto.BlogResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.BlogResponse, error)
	Update(id uint, req dto.UpdateBlogRequest) (*dto.BlogResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, error)
}

type blogService struct {
	repo  repository.BlogRepository
	cache cache.CacheService
}

func NewBlogService(repo repository.BlogRepository, cache cache.CacheService) BlogService {
	return &blogService{
		repo:  repo,
		cache: cache,
	}
}

func (s *blogService) Create(req dto.CreateBlogRequest) (*dto.BlogResponse, error) {
	entity := &entity.Blog{
		Name:      req.Name,
		Content:   req.Content,
		Thumbnail: req.Thumbnail,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("blog:*")

	logger.AuditLogger.Info().
		Uint("blog_id", entity.ID).
		Str("action", "blog_creation").
		Msg("blog created")

	return s.mapToResponse(entity), nil
}

func (s *blogService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("blog:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.BlogResponse
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

func (s *blogService) GetByID(id uint) (*dto.BlogResponse, error) {
	cacheKey := fmt.Sprintf("blog:%d", id)
	var cached dto.BlogResponse
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

func (s *blogService) Update(id uint, req dto.UpdateBlogRequest) (*dto.BlogResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	entity.Content = req.Content
	if req.Thumbnail != "" {
		entity.Thumbnail = req.Thumbnail
	}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("blog:%d", id))
	s.cache.DeletePattern("blog:*")

	logger.AuditLogger.Info().
		Uint("blog_id", entity.ID).
		Str("action", "blog_update").
		Msg("blog updated")

	return s.mapToResponse(entity), nil
}

func (s *blogService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("blog:%d", id))
	s.cache.DeletePattern("blog:*")

	logger.AuditLogger.Info().
		Uint("blog_id", id).
		Str("action", "blog_deletion").
		Msg("blog deleted")

	return s.repo.Delete(id)
}

func (s *blogService) Export(format string) ([]byte, error) {
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

func (s *blogService) generateCSV(entities []entity.Blog) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "Name", "Content", "Thumbnail", "Created At"}
	writer.Write(header)

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			fmt.Sprintf("%v", e.Name),
			fmt.Sprintf("%v", e.Content),
			fmt.Sprintf("%v", e.Thumbnail),
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *blogService) generateExcel(entities []entity.Blog) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	header := []string{"ID", "Name", "Content", "Thumbnail", "Created At"}
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
		f.SetCellValue(sheet, cell, e.Content)
		cell, _ = excelize.CoordinatesToCellName(2+2, rowNum)
		f.SetCellValue(sheet, cell, e.Thumbnail)
		lastCell, _ := excelize.CoordinatesToCellName(len(header), rowNum)
		f.SetCellValue(sheet, lastCell, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *blogService) mapToResponse(entity *entity.Blog) *dto.BlogResponse {
	return &dto.BlogResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		Content:   entity.Content,
		Thumbnail: entity.Thumbnail,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
