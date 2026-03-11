package service

import (
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type TestsajaService interface {
	Create(req dto.CreateTestsajaRequest) (*dto.TestsajaResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.TestsajaResponse, error)
	Update(id uint, req dto.UpdateTestsajaRequest) (*dto.TestsajaResponse, error)
	Delete(id uint) error
}

type testsajaService struct {
	repo  repository.TestsajaRepository
	cache cache.CacheService
}

func NewTestsajaService(repo repository.TestsajaRepository, cache cache.CacheService) TestsajaService {
	return &testsajaService{
		repo:  repo,
		cache: cache,
	}
}

func (s *testsajaService) Create(req dto.CreateTestsajaRequest) (*dto.TestsajaResponse, error) {
	entity := &entity.Testsaja{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("testsaja:*")

	logger.AuditLogger.Info().
		Uint("testsaja_id", entity.ID).
		Str("action", "testsaja_creation").
		Msg("testsaja created")

	return s.mapToResponse(entity), nil
}

func (s *testsajaService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("testsaja:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.TestsajaResponse
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

func (s *testsajaService) GetByID(id uint) (*dto.TestsajaResponse, error) {
	cacheKey := fmt.Sprintf("testsaja:%d", id)
	var cached dto.TestsajaResponse
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

func (s *testsajaService) Update(id uint, req dto.UpdateTestsajaRequest) (*dto.TestsajaResponse, error) {
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

	s.cache.Delete(fmt.Sprintf("testsaja:%d", id))
	s.cache.DeletePattern("testsaja:*")

	logger.AuditLogger.Info().
		Uint("testsaja_id", entity.ID).
		Str("action", "testsaja_update").
		Msg("testsaja updated")

	return s.mapToResponse(entity), nil
}

func (s *testsajaService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("testsaja:%d", id))
	s.cache.DeletePattern("testsaja:*")

	logger.AuditLogger.Info().
		Uint("testsaja_id", id).
		Str("action", "testsaja_deletion").
		Msg("testsaja deleted")

	return s.repo.Delete(id)
}

func (s *testsajaService) mapToResponse(entity *entity.Testsaja) *dto.TestsajaResponse {
	return &dto.TestsajaResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
