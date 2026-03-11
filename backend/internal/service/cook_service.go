package service

import (
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	defaultdto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type CookService interface {
	Create(req dto.CreateCookRequest) (*dto.CookResponse, error)
	GetAll(pagination *defaultdto.PaginationRequest) (*defaultdto.PaginationResponse, error)
	GetByID(id uint) (*dto.CookResponse, error)
	Update(id uint, req dto.UpdateCookRequest) (*dto.CookResponse, error)
	Delete(id uint) error
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

func (s *cookService) mapToResponse(entity *entity.Cook) *dto.CookResponse {
	return &dto.CookResponse{
		ID:        entity.ID,
		Name: entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
