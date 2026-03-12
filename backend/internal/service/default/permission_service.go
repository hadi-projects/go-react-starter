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

type PermissionService interface {
	Create(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdatePermissionRequest) (*dto.PermissionResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, format string) ([]byte, string, error)
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

func (s *permissionService) Create(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	permission := &entity.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}

	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	// logger.AuditLogger.Info().
	// 	Uint("permission_id", permission.ID).
	// 	Str("name", permission.Name).
	// 	Str("action", "permission_creation").
	// 	Msg("permission created")
	logger.LogAudit(ctx, "CREATE", "PERMISSION", fmt.Sprintf("%d", permission.ID), fmt.Sprintf("name: %s", permission.Name))

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

func (s *permissionService) Update(ctx context.Context, id uint, req dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
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

	// logger.AuditLogger.Info().
	// 	Uint("permission_id", permission.ID).
	// 	Str("name", permission.Name).
	// 	Str("action", "permission_update").
	// 	Msg("permission updated")
	logger.LogAudit(ctx, "UPDATE", "PERMISSION", fmt.Sprintf("%d", id), fmt.Sprintf("name: %s", permission.Name))

	return &dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

func (s *permissionService) Delete(ctx context.Context, id uint) error {
	// Invalidate permissions list cache
	s.cache.DeletePattern("permissions:*")

	// logger.AuditLogger.Info().
	// 	Uint("target_permission_id", id).
	// 	Str("action", "permission_deletion").
	// 	Msg("permission deleted")
	logger.LogAudit(ctx, "DELETE", "PERMISSION", fmt.Sprintf("%d", id), "")

	return s.repo.Delete(id)
}

func (s *permissionService) Export(ctx context.Context, format string) ([]byte, string, error) {
	pagination := &dto.PaginationRequest{
		Page:  1,
		Limit: 1000000,
	}

	permissions, _, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(permissions)
	}
	return s.generateExcel(permissions)
}

func (s *permissionService) generateCSV(permissions []entity.Permission) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Name", "Description", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, perm := range permissions {
		row := []string{
			fmt.Sprintf("%d", perm.ID),
			perm.Name,
			perm.Description,
			perm.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "permissions.csv", nil
}

func (s *permissionService) generateExcel(permissions []entity.Permission) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Permissions"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Name", "Description", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, perm := range permissions {
		row := []interface{}{
			perm.ID,
			perm.Name,
			perm.Description,
			perm.CreatedAt.Format("2006-01-02 15:04:05"),
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

	return buf.Bytes(), "permissions.xlsx", nil
}
