package service

import (
	"fmt"
	"math"
	"time"

	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type ProdukService interface {
	Create(req dto.CreateProdukRequest) (*dto.ProdukResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.ProdukResponse, error)
	Update(id uint, req dto.UpdateProdukRequest) (*dto.ProdukResponse, error)
	Delete(id uint) error
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

func (s *produkService) Create(req dto.CreateProdukRequest) (*dto.ProdukResponse, error) {
	entity := &entity.Produk{
		Name: req.Name,
		Harga: req.Harga,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("produk:*")

	
	logger.AuditLogger.Info().
		Uint("produk_id", entity.ID).
		Str("action", "produk_creation").
		Msg("produk created")
	

	return s.mapToResponse(entity), nil
}

func (s *produkService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("produk:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
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

	s.cache.Set(cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *produkService) GetByID(id uint) (*dto.ProdukResponse, error) {
	cacheKey := fmt.Sprintf("produk:%d", id)
	var cached dto.ProdukResponse
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

func (s *produkService) Update(id uint, req dto.UpdateProdukRequest) (*dto.ProdukResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	entity.Harga = req.Harga

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("produk:%d", id))
	s.cache.DeletePattern("produk:*")

	
	logger.AuditLogger.Info().
		Uint("produk_id", entity.ID).
		Str("action", "produk_update").
		Msg("produk updated")
	

	return s.mapToResponse(entity), nil
}

func (s *produkService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("produk:%d", id))
	s.cache.DeletePattern("produk:*")

	
	logger.AuditLogger.Info().
		Uint("produk_id", id).
		Str("action", "produk_deletion").
		Msg("produk deleted")
	

	return s.repo.Delete(id)
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
