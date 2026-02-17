package service

import (
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type AbcService interface {
	Create(req dto.CreateAbcRequest) (*dto.AbcResponse, error)
	GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	GetByID(id uint) (*dto.AbcResponse, error)
	Update(id uint, req dto.UpdateAbcRequest) (*dto.AbcResponse, error)
	Delete(id uint) error
}

type abcService struct {
	repo  repository.AbcRepository
	cache cache.CacheService
}

func NewAbcService(repo repository.AbcRepository, cache cache.CacheService) AbcService {
	return &abcService{
		repo:  repo,
		cache: cache,
	}
}

func (s *abcService) Create(req dto.CreateAbcRequest) (*dto.AbcResponse, error) {
	entity := &entity.Abc{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("abc:*")

	logger.AuditLogger.Info().
		Uint("abc_id", entity.ID).
		Str("action", "abc_creation").
		Msg("abc created")

	return s.mapToResponse(entity), nil
}

func (s *abcService) GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("abc:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached dto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.AbcResponse
	for _, e := range entities {
		responses = append(responses, *s.mapToResponse(&e))
	}

	response := &dto.PaginationResponse{
		Data: responses,
		Meta: dto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  int(math.Ceil(float64(total) / float64(pagination.GetLimit()))),
			TotalItems:  total,
			Limit:       pagination.GetLimit(),
		},
	}

	s.cache.Set(cacheKey, response, 5*time.Minute)
	return response, nil
}

func (s *abcService) GetByID(id uint) (*dto.AbcResponse, error) {
	cacheKey := fmt.Sprintf("abc:%d", id)
	var cached dto.AbcResponse
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

func (s *abcService) Update(id uint, req dto.UpdateAbcRequest) (*dto.AbcResponse, error) {
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

	s.cache.Delete(fmt.Sprintf("abc:%d", id))
	s.cache.DeletePattern("abc:*")

	logger.AuditLogger.Info().
		Uint("abc_id", entity.ID).
		Str("action", "abc_update").
		Msg("abc updated")

	return s.mapToResponse(entity), nil
}

func (s *abcService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("abc:%d", id))
	s.cache.DeletePattern("abc:*")

	logger.AuditLogger.Info().
		Uint("abc_id", id).
		Str("action", "abc_deletion").
		Msg("abc deleted")

	return s.repo.Delete(id)
}

func (s *abcService) mapToResponse(entity *entity.Abc) *dto.AbcResponse {
	return &dto.AbcResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
