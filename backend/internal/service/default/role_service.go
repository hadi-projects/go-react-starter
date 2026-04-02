package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type RoleService interface {
	Create(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetAll(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.RoleResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, format string) ([]byte, string, error)
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

func (s *roleService) Create(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &entity.Role{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	}

	if err := s.roleRepo.Create(ctx, role, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Invalidate roles list cache
	s.cache.DeletePattern(ctx, "roles:*")

	// logger.AuditLogger.Info().
	// 	Uint("role_id", role.ID).
	// 	Str("name", role.Name).
	// 	Str("action", "role_creation").
	// 	Msg("role created")
	logger.LogAudit(ctx, "CREATE", "ROLE", fmt.Sprintf("%d", role.ID), fmt.Sprintf("name: %s", role.Name))

	// Fetch again to get permissions populated (or we can construct response manually if we trust repo)
	// Better to fetch to be sure.
	createdRole, err := s.roleRepo.FindByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(createdRole), nil
}

func (s *roleService) GetAll(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("roles:page:%d:limit:%d:search:%s:cat:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search, pagination.Category)
	var cached dto.PaginationResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	roles, total, err := s.roleRepo.FindAll(ctx, pagination)
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
	s.cache.Set(ctx, cacheKey, response, ttl)

	return response, nil
}

func (s *roleService) GetByID(ctx context.Context, id uint) (*dto.RoleResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("role:%d", id)
	var cached dto.RoleResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.mapToResponse(role)

	// Cache the result
	ttl := time.Duration(300) * time.Second
	s.cache.Set(ctx, cacheKey, response, ttl)

	return response, nil
}

func (s *roleService) Update(ctx context.Context, id uint, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if req.Category != "" {
		role.Category = req.Category
	}

	if err := s.roleRepo.Update(ctx, role, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Invalidate role cache, roles list cache, AND all user caches since permissions changed
	s.cache.Delete(ctx, fmt.Sprintf("role:%d", id))
	s.cache.DeletePattern(ctx, "roles:*")
	s.cache.DeletePattern(ctx, "user:*")

	// logger.AuditLogger.Info().
	// 	Uint("role_id", role.ID).
	// 	Str("name", role.Name).
	// 	Str("action", "role_update").
	// 	Msg("role updated")
	logger.LogAudit(ctx, "UPDATE", "ROLE", fmt.Sprintf("%d", id), fmt.Sprintf("name: %s", role.Name))

	updatedRole, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(updatedRole), nil
}

func (s *roleService) Delete(ctx context.Context, id uint) error {
	// Invalidate role cache, roles list cache, AND all user caches
	s.cache.Delete(ctx, fmt.Sprintf("role:%d", id))
	s.cache.DeletePattern(ctx, "roles:*")
	s.cache.DeletePattern(ctx, "user:*")

	// logger.AuditLogger.Info().
	// 	Uint("target_role_id", id).
	// 	Str("action", "role_deletion").
	// 	Msg("role deleted")
	logger.LogAudit(ctx, "DELETE", "ROLE", fmt.Sprintf("%d", id), "")

	return s.roleRepo.Delete(ctx, id)
}

func (s *roleService) Export(ctx context.Context, format string) ([]byte, string, error) {
	pagination := &dto.PaginationRequest{
		Page:  1,
		Limit: 1000000,
	}

	roles, _, err := s.roleRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(roles)
	}
	return s.generateExcel(roles)
}

func (s *roleService) generateCSV(roles []entity.Role) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Name", "Description", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, role := range roles {
		row := []string{
			fmt.Sprintf("%d", role.ID),
			role.Name,
			role.Description,
			role.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "roles.csv", nil
}

func (s *roleService) generateExcel(roles []entity.Role) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Roles"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Name", "Description", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, role := range roles {
		row := []interface{}{
			role.ID,
			role.Name,
			role.Description,
			role.CreatedAt.Format("2006-01-02 15:04:05"),
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

	return buf.Bytes(), "roles.xlsx", nil
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
		Category:    role.Category,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
