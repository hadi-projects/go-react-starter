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

type ProdukService interface {
	Create(ctx context.Context, req dto.CreateProdukRequest) (*dto.ProdukResponse, error)
	GetAll(ctx context.Context, pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.ProdukResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateProdukRequest) (*dto.ProdukResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, format string) ([]byte, string, error)
}

type produkService struct {
	repo  repository.ProdukRepository
	cache cache.CacheService
}

func NewProdukService(repo repository.ProdukRepository, cache cache.CacheService) ProdukService {
	return &produkService{
		repo:  repo,
		cache: cache,
	}
}

func (s *produkService) Create(ctx context.Context, req dto.CreateProdukRequest) (*dto.ProdukResponse, error) {
	entity := &entity.Produk{
		Name: req.Name,
		Harga: req.Harga,
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern(ctx, "produk:*")

	
	// logger.AuditLogger.Info().
	// 	Uint("produk_id", entity.ID).
	// 	Str("action", "produk_creation").
	// 	Msg("produk created")
	logger.LogAudit(ctx, "CREATE", "PRODUK", fmt.Sprintf("%d", entity.ID), fmt.Sprintf("name: %s, harga: %d", entity.Name, entity.Harga))
	

	return s.mapToResponse(entity), nil
}

func (s *produkService) GetAll(ctx context.Context, pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("produk:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.ProdukResponse
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

	s.cache.Set(ctx, cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *produkService) GetByID(ctx context.Context, id uint) (*dto.ProdukResponse, error) {
	cacheKey := fmt.Sprintf("produk:%d", id)
	var cached dto.ProdukResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.mapToResponse(entity)
	s.cache.Set(ctx, cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *produkService) Update(ctx context.Context, id uint, req dto.UpdateProdukRequest) (*dto.ProdukResponse, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	entity.Harga = req.Harga

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	s.cache.Delete(ctx, fmt.Sprintf("produk:%d", id))
	s.cache.DeletePattern(ctx, "produk:*")

	
	// logger.AuditLogger.Info().
	// 	Uint("produk_id", entity.ID).
	// 	Str("action", "produk_update").
	// 	Msg("produk updated")
	logger.LogAudit(ctx, "UPDATE", "PRODUK", fmt.Sprintf("%d", id), fmt.Sprintf("name: %s, harga: %d", entity.Name, entity.Harga))
	

	return s.mapToResponse(entity), nil
}

func (s *produkService) Delete(ctx context.Context, id uint) error {
	s.cache.Delete(ctx, fmt.Sprintf("produk:%d", id))
	s.cache.DeletePattern(ctx, "produk:*")

	
	// logger.AuditLogger.Info().
	// 	Uint("produk_id", id).
	// 	Str("action", "produk_deletion").
	// 	Msg("produk deleted")
	logger.LogAudit(ctx, "DELETE", "PRODUK", fmt.Sprintf("%d", id), "")
	

	return s.repo.Delete(ctx, id)
}

func (s *produkService) Export(ctx context.Context, format string) ([]byte, string, error) {
	pagination := &defaultDto.PaginationRequest{
		Page:  1,
		Limit: 1000000,
	}

	entities, _, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(entities)
	}
	return s.generateExcel(entities)
}

func (s *produkService) generateCSV(entities []entity.Produk) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Name", "Harga", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			e.Name,
			fmt.Sprintf("%d", e.Harga),
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "produk.csv", nil
}

func (s *produkService) generateExcel(entities []entity.Produk) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Produk"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Name", "Harga", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, e := range entities {
		row := []interface{}{
			e.ID,
			e.Name,
			e.Harga,
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

	return buf.Bytes(), "produk.xlsx", nil
}

func (s *produkService) mapToResponse(entity *entity.Produk) *dto.ProdukResponse {
	return &dto.ProdukResponse{
		ID:        entity.ID,
		Name: entity.Name,
		Harga: entity.Harga,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
