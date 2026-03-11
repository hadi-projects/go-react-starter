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

type RoleService interface {
	Create(req dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	GetByID(id uint) (*dto.RoleResponse, error)
	Update(id uint, req dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	Delete(id uint) error
}

type roleService struct {
	roleRepo repository.RoleRepository
	cache    cache.CacheService
}

func NewRoleService(roleRepo repository.RoleRepository, cache cache.CacheService) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		cache:    cache,
	}
}

func (s *roleService) Create(req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &entity.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(role, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Invalidate roles list cache
	s.cache.DeletePattern("roles:*")

	logger.AuditLogger.Info().
		Uint("role_id", role.ID).
		Str("name", role.Name).
		Str("action", "role_creation").
		Msg("role created")

	// Fetch again to get permissions populated (or we can construct response manually if we trust repo)
	// Better to fetch to be sure.
	createdRole, err := s.roleRepo.FindByID(role.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(createdRole), nil
}

func (s *roleService) GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("roles:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached dto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	roles, total, err := s.roleRepo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.mapToResponse(&role))
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
	ttl := time.Duration(300) * time.Second // Default 5 minutes
	s.cache.Set(cacheKey, response, ttl)

	return response, nil
}

func (s *roleService) GetByID(id uint) (*dto.RoleResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("role:%d", id)
	var cached dto.RoleResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.mapToResponse(role)

	// Cache the result
	ttl := time.Duration(300) * time.Second
	s.cache.Set(cacheKey, response, ttl)

	return response, nil
}

func (s *roleService) Update(id uint, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(role, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Invalidate role cache, roles list cache, AND all user caches since permissions changed
	s.cache.Delete(fmt.Sprintf("role:%d", id))
	s.cache.DeletePattern("roles:*")
	s.cache.DeletePattern("user:*")

	logger.AuditLogger.Info().
		Uint("role_id", role.ID).
		Str("name", role.Name).
		Str("action", "role_update").
		Msg("role updated")

	updatedRole, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(updatedRole), nil
}

func (s *roleService) Delete(id uint) error {
	// Invalidate role cache, roles list cache, AND all user caches
	s.cache.Delete(fmt.Sprintf("role:%d", id))
	s.cache.DeletePattern("roles:*")
	s.cache.DeletePattern("user:*")

	logger.AuditLogger.Info().
		Uint("target_role_id", id).
		Str("action", "role_deletion").
		Msg("role deleted")

	return s.roleRepo.Delete(id)
}

func (s *roleService) mapToResponse(role *entity.Role) *dto.RoleResponse {
	var permissions []dto.PermissionResponse
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.PermissionResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
