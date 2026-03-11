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

type TestduaService interface {
	Create(req dto.CreateTestduaRequest) (*dto.TestduaResponse, error)
	GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetByID(id uint) (*dto.TestduaResponse, error)
	Update(id uint, req dto.UpdateTestduaRequest) (*dto.TestduaResponse, error)
	Delete(id uint) error
}

type testduaService struct {
	repo  repository.TestduaRepository
	cache cache.CacheService
}

func NewTestduaService(repo repository.TestduaRepository, cache cache.CacheService) TestduaService {
	return &testduaService{
		repo:  repo,
		cache: cache,
	}
}

func (s *testduaService) Create(req dto.CreateTestduaRequest) (*dto.TestduaResponse, error) {
	entity := &entity.Testdua{
		Name: req.Name,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("testdua:*")

	
	logger.AuditLogger.Info().
		Uint("testdua_id", entity.ID).
		Str("action", "testdua_creation").
		Msg("testdua created")
	

	return s.mapToResponse(entity), nil
}

func (s *testduaService) GetAll(pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("testdua:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.TestduaResponse
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

func (s *testduaService) GetByID(id uint) (*dto.TestduaResponse, error) {
	cacheKey := fmt.Sprintf("testdua:%d", id)
	var cached dto.TestduaResponse
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

func (s *testduaService) Update(id uint, req dto.UpdateTestduaRequest) (*dto.TestduaResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	entity.Name = req.Name

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("testdua:%d", id))
	s.cache.DeletePattern("testdua:*")

	
	logger.AuditLogger.Info().
		Uint("testdua_id", entity.ID).
		Str("action", "testdua_update").
		Msg("testdua updated")
	

	return s.mapToResponse(entity), nil
}

func (s *testduaService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("testdua:%d", id))
	s.cache.DeletePattern("testdua:*")

	
	logger.AuditLogger.Info().
		Uint("testdua_id", id).
		Str("action", "testdua_deletion").
		Msg("testdua deleted")
	

	return s.repo.Delete(id)
}

func (s *testduaService) mapToResponse(entity *entity.Testdua) *dto.TestduaResponse {
	return &dto.TestduaResponse{
		ID:        entity.ID,
		Name: entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
