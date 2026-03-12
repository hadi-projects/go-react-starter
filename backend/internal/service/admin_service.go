package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/internal/dto"
	defaultdto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type AdminService interface {
	Create(req dto.CreateAdminRequest) (*dto.AdminResponse, error)
	GetAll(pagination *defaultdto.PaginationRequest) (*defaultdto.PaginationResponse, error)
	GetByID(id uint) (*dto.AdminResponse, error)
	Update(id uint, req dto.UpdateAdminRequest) (*dto.AdminResponse, error)
	Delete(id uint) error
	Export(format string) ([]byte, string, error)
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

func (s *adminService) Export(format string) ([]byte, string, error) {
	pagination := &defaultdto.PaginationRequest{
		Page:  1,
		Limit: 1000000,
	}

	entities, _, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(entities)
	}
	return s.generateExcel(entities)
}

func (s *adminService) generateCSV(entities []entity.Admin) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"ID", "Name", "Email", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, e := range entities {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			e.Name,
			e.Email,
			e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "admin.csv", nil
}

func (s *adminService) generateExcel(entities []entity.Admin) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Admin"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Name", "Email", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, e := range entities {
		row := []interface{}{
			e.ID,
			e.Name,
			e.Email,
			e.CreatedAt.Format("2006-01-02 15:04:05"),
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

	return buf.Bytes(), "admin.xlsx", nil
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
