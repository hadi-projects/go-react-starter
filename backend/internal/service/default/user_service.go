package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/config"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetMe(ctx context.Context, userID uint) (*dto.AuthUserResponse, error)
	GetAll(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, format string) ([]byte, string, error)
}

type userService struct {
	userRepo repository.UserRepository
	config   *config.Config
	cache    cache.CacheService
}

func NewUserService(userRepo repository.UserRepository, config *config.Config, cache cache.CacheService) UserService {
	return &userService{
		userRepo: userRepo,
		config:   config,
		cache:    cache,
	}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.Security.BCryptCost)
	if err != nil {
		return nil, err
	}

	roleID := uint(2) // Default fallback
	role, err := s.userRepo.FindRoleByName(ctx, "user")
	if err == nil {
		roleID = role.ID
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   roleID,
		Status:   "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Invalidate users list cache
	s.cache.DeletePattern(ctx, "users:*")

	// logger.AuditLogger.Info().
	// 	Uint("user_id", user.ID).
	// 	Str("email", user.Email).
	// 	Str("action", "user_registration").
	// 	Msg("user registered successfully")
	logger.LogAudit(ctx, "REGISTER", "USER", fmt.Sprintf("%d", user.ID), fmt.Sprintf("email: %s", user.Email))

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.Security.BCryptCost)
	if err != nil {
		return nil, err
	}

	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   req.RoleID,
		Status:   status,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Invalidate users list cache
	s.cache.DeletePattern(ctx, "users:*")

	// logger.AuditLogger.Info().
	// 	Uint("user_id", user.ID).
	// 	Str("email", user.Email).
	// 	Str("action", "user_creation").
	// 	Msg("user created by admin")
	logger.LogAudit(ctx, "CREATE", "USER", fmt.Sprintf("%d", user.ID), fmt.Sprintf("email: %s, role_id: %d", user.Email, user.RoleID))

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetMe(ctx context.Context, userID uint) (*dto.AuthUserResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:%d", userID)
	var cached dto.AuthUserResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var permissionsMask uint64
	if user.RoleID != 0 {
		for _, p := range user.Role.Permissions {
			if p.ID <= 64 {
				permissionsMask |= (1 << (p.ID - 1))
			}
		}
	}

	response := &dto.AuthUserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		RoleID:          user.RoleID,
		Role:            user.Role.Name,
		PermissionsMask: permissionsMask,
		Status:          user.Status,
	}

	// Cache the result
	ttl := time.Duration(s.config.Redis.TTL) * time.Second
	s.cache.Set(ctx, cacheKey, response, ttl)

	return response, nil
}

func (s *userService) GetAll(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("users:page:%d:limit:%d:search:%s", pagination.GetPage(), pagination.GetLimit(), pagination.Search)
	var cached dto.PaginationResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	users, total, err := s.userRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			RoleID:    user.RoleID,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response := &dto.PaginationResponse{
		Data: userResponses,
		Meta: dto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  int(math.Ceil(float64(total) / float64(pagination.GetLimit()))),
			TotalData:   total,
			Limit:       pagination.GetLimit(),
		},
	}

	// Cache the result
	ttl := time.Duration(s.config.Redis.TTL) * time.Second
	s.cache.Set(ctx, cacheKey, response, ttl)

	return response, nil
}

func (s *userService) Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.Security.BCryptCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Invalidate user cache and users list cache
	s.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
	s.cache.DeletePattern(ctx, "users:*")

	// logger.AuditLogger.Info().
	// 	Uint("user_id", user.ID).
	// 	Str("email", user.Email).
	// 	Str("action", "user_update").
	// 	Msg("user details updated")
	logger.LogAudit(ctx, "UPDATE", "USER", fmt.Sprintf("%d", id), fmt.Sprintf("email: %s, role_id: %d", user.Email, user.RoleID))

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	// Invalidate user cache and users list cache
	s.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
	s.cache.DeletePattern(ctx, "users:*")

	// logger.AuditLogger.Info().
	// 	Uint("target_user_id", id).
	// 	Str("action", "user_deletion").
	// 	Msg("user deleted")
	logger.LogAudit(ctx, "DELETE", "USER", fmt.Sprintf("%d", id), "")

	return s.userRepo.Delete(ctx, id)
}

func (s *userService) Export(ctx context.Context, format string) ([]byte, string, error) {
	pagination := &dto.PaginationRequest{
		Page:  1,
		Limit: 1000000, // Export all users
	}

	users, _, err := s.userRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, "", err
	}

	if format == "csv" {
		return s.generateCSV(users)
	}
	return s.generateExcel(users)
}

func (s *userService) generateCSV(users []entity.User) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Header
	header := []string{"ID", "Name", "Email", "Role", "Status", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Data
	for _, user := range users {
		row := []string{
			fmt.Sprintf("%d", user.ID),
			user.Name,
			user.Email,
			user.Role.Name,
			user.Status,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	return buf.Bytes(), "users.csv", nil
}

func (s *userService) generateExcel(users []entity.User) ([]byte, string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := "Users"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, "", err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Header
	headers := []string{"ID", "Name", "Email", "Role", "Status", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data
	for i, user := range users {
		row := []interface{}{
			user.ID,
			user.Name,
			user.Email,
			user.Role.Name,
			user.Status,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
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

	return buf.Bytes(), "users.xlsx", nil
}
