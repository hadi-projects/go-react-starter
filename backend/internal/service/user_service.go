package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	GetMe(userID uint) (*dto.UserResponse, error)
	GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error)
	Update(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id uint) error
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

func (s *userService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.Security.BCryptCost)
	if err != nil {
		return nil, err
	}

	roleID := uint(2) // Default fallback
	role, err := s.userRepo.FindRoleByName("user")
	if err == nil {
		roleID = role.ID
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   roleID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Invalidate users list cache
	s.cache.DeletePattern("users:*")

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetMe(userID uint) (*dto.UserResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:%d", userID)
	var cached dto.UserResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Cache the result
	ttl := time.Duration(s.config.Redis.TTL) * time.Second
	s.cache.Set(cacheKey, response, ttl)

	return response, nil
}

func (s *userService) GetAll(pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("users:page:%d:limit:%d", pagination.GetPage(), pagination.GetLimit())
	var cached dto.PaginationResponse
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	users, total, err := s.userRepo.FindAll(pagination)
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
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response := &dto.PaginationResponse{
		Data: userResponses,
		Meta: dto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  int(math.Ceil(float64(total) / float64(pagination.GetLimit()))),
			TotalItems:  total,
			Limit:       pagination.GetLimit(),
		},
	}

	// Cache the result
	ttl := time.Duration(s.config.Redis.TTL) * time.Second
	s.cache.Set(cacheKey, response, ttl)

	return response, nil
}

func (s *userService) Update(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
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

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Invalidate user cache and users list cache
	s.cache.Delete(fmt.Sprintf("user:%d", id))
	s.cache.DeletePattern("users:*")

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) Delete(id uint) error {
	// Invalidate user cache and users list cache
	s.cache.Delete(fmt.Sprintf("user:%d", id))
	s.cache.DeletePattern("users:*")

	return s.userRepo.Delete(id)
}
