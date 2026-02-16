package service

import (
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
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
		Name: req.Name,
	}

	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}

	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	return &dto.PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}, nil
}

func (s *permissionService) GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("permissions:page:%d:limit:%d", pagination.GetPage(), pagination.GetLimit())
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
			ID:        perm.ID,
			Name:      perm.Name,
			CreatedAt: perm.CreatedAt,
			UpdatedAt: perm.UpdatedAt,
		})
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
	if err := s.repo.Update(permission); err != nil {
		return nil, err
	}

	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	return &dto.PermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}, nil
}

func (s *permissionService) Delete(id uint) error {
	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	return s.repo.Delete(id)
}
