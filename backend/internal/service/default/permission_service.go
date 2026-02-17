package service

import (
	"fmt"
	"math"
	"time"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type PermissionService interface {
	Create(req dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	Update(id uint, req dto.UpdatePermissionRequest) (*dto.PermissionResponse, error)
	Delete(id uint) error
}

type permissionService struct {
	repo  repository.PermissionRepository
	cache cache.CacheService
}

func NewPermissionService(repo repository.PermissionRepository, cache cache.CacheService) PermissionService {
	return &permissionService{
		repo:  repo,
		cache: cache,
	}
}

func (s *permissionService) Create(req dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	permission := &entity.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}

	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	logger.AuditLogger.Info().
		Uint("permission_id", permission.ID).
		Str("name", permission.Name).
		Str("action", "permission_creation").
		Msg("permission created")

	return &dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

func (s *permissionService) GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("permissions:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached dto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	permissions, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.PermissionResponse
	for _, perm := range permissions {
		responses = append(responses, dto.PermissionResponse{
			ID:          perm.ID,
			Name:        perm.Name,
			Description: perm.Description,
			CreatedAt:   perm.CreatedAt,
			UpdatedAt:   perm.UpdatedAt,
		})
	}

	response := &dto.PaginationResponse{
		Data: responses,
		Meta: dto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  int(math.Ceil(float64(total) / float64(pagination.GetLimit()))),
			TotalData:   total,
			Limit:       pagination.GetLimit(),
		},
	}

	// Cache the result
	ttl := time.Duration(300) * time.Second
	s.cache.Set(cacheKey, response, ttl)

	return response, nil
}

func (s *permissionService) Update(id uint, req dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	permission.Name = req.Name
	permission.Description = req.Description
	if err := s.repo.Update(permission); err != nil {
		return nil, err
	}

	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	logger.AuditLogger.Info().
		Uint("permission_id", permission.ID).
		Str("name", permission.Name).
		Str("action", "permission_update").
		Msg("permission updated")

	return &dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

func (s *permissionService) Delete(id uint) error {
	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	logger.AuditLogger.Info().
		Uint("target_permission_id", id).
		Str("action", "permission_deletion").
		Msg("permission deleted")

	return s.repo.Delete(id)
}
