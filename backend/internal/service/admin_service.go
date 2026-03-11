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

type AdminService interface {
	Create(req dto.CreateAdminRequest) (*dto.AdminResponse, error)
	GetAll(pagination *defaultdto.PaginationRequest) (*defaultdto.PaginationResponse, error)
	GetByID(id uint) (*dto.AdminResponse, error)
	Update(id uint, req dto.UpdateAdminRequest) (*dto.AdminResponse, error)
	Delete(id uint) error
}

type adminService struct {
	repo  repository.AdminRepository
	cache cache.CacheService
}

func NewAdminService(repo repository.AdminRepository, cache cache.CacheService) AdminService {
	return &adminService{
		repo:  repo,
		cache: cache,
	}
}

func (s *adminService) Create(req dto.CreateAdminRequest) (*dto.AdminResponse, error) {
	entity := &entity.Admin{
		Name: req.Name,
		Email: req.Email,
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	s.cache.DeletePattern("admin:*")

	
	logger.AuditLogger.Info().
		Uint("admin_id", entity.ID).
		Str("action", "admin_creation").
		Msg("admin created")
	

	return s.mapToResponse(entity), nil
}

func (s *adminService) GetAll(pagination *defaultdto.PaginationRequest) (*defaultdto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("admin:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached defaultdto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	entities, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.AdminResponse
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

func (s *adminService) GetByID(id uint) (*dto.AdminResponse, error) {
	cacheKey := fmt.Sprintf("admin:%d", id)
	var cached dto.AdminResponse
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

func (s *adminService) Update(id uint, req dto.UpdateAdminRequest) (*dto.AdminResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Email != "" {
		entity.Email = req.Email
	}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	s.cache.Delete(fmt.Sprintf("admin:%d", id))
	s.cache.DeletePattern("admin:*")

	
	logger.AuditLogger.Info().
		Uint("admin_id", entity.ID).
		Str("action", "admin_update").
		Msg("admin updated")
	

	return s.mapToResponse(entity), nil
}

func (s *adminService) Delete(id uint) error {
	s.cache.Delete(fmt.Sprintf("admin:%d", id))
	s.cache.DeletePattern("admin:*")

	
	logger.AuditLogger.Info().
		Uint("admin_id", id).
		Str("action", "admin_deletion").
		Msg("admin deleted")
	

	return s.repo.Delete(id)
}

func (s *adminService) mapToResponse(entity *entity.Admin) *dto.AdminResponse {
	return &dto.AdminResponse{
		ID:        entity.ID,
		Name: entity.Name,
		Email: entity.Email,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
