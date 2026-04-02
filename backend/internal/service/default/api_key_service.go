package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
)

type ApiKeyService interface {
	Generate(ctx context.Context, userID uint, req dto.ApiKeyCreateRequest) (*dto.ApiKeyCreateResponse, error)
	GetAll(ctx context.Context, pagination *dto.PaginationRequest) ([]dto.ApiKeyResponse, int64, error)
	Delete(ctx context.Context, id uint) error
	Validate(ctx context.Context, rawKey string) (*entity.ApiKey, uint64, error)
}

type apiKeyService struct {
	repo     repository.ApiKeyRepository
	roleRepo repository.RoleRepository
	cache    cache.CacheService
}

func NewApiKeyService(repo repository.ApiKeyRepository, roleRepo repository.RoleRepository, cache cache.CacheService) ApiKeyService {
	return &apiKeyService{repo: repo, roleRepo: roleRepo, cache: cache}
}

func (s *apiKeyService) Generate(ctx context.Context, userID uint, req dto.ApiKeyCreateRequest) (*dto.ApiKeyCreateResponse, error) {
	var rawKey string
	var prefix string

	if req.Type == "uuid" {
		rawKey = uuid.New().String()
		prefix = rawKey[:8]
	} else {
		// Private key format: sk_tp_ + random
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		randomStr := base64.URLEncoding.EncodeToString(b)
		rawKey = "sk_tp_" + randomStr
		prefix = "sk_tp_..." + randomStr[len(randomStr)-4:]
	}

	hash := s.hashKey(rawKey)

	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		exp := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &exp
	}

	role, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("role not found")
	}
	if role.Category != "api" {
		return nil, fmt.Errorf("role category must be 'api' for API Keys")
	}

	apiKey := &entity.ApiKey{
		UserID:     userID,
		Name:       req.Name,
		KeyHash:    hash,
		Prefix:     prefix,
		RoleID:     req.RoleID,
		AllowedIPs: req.AllowedIPs,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &dto.ApiKeyCreateResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		RawKey:    rawKey, // Only shown once
		Prefix:    apiKey.Prefix,
		RoleName:  "...", // Handled by caller or through reload
		ExpiresAt: apiKey.ExpiresAt,
	}, nil
}

func (s *apiKeyService) GetAll(ctx context.Context, pagination *dto.PaginationRequest) ([]dto.ApiKeyResponse, int64, error) {
	keys, total, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.ApiKeyResponse
	for _, k := range keys {
		res = append(res, dto.ApiKeyResponse{
			ID:         k.ID,
			Name:       k.Name,
			Prefix:     k.Prefix,
			RoleName:   k.Role.Name,
			AllowedIPs: k.AllowedIPs,
			ExpiresAt:  k.ExpiresAt,
			LastUsedAt: k.LastUsedAt,
			CreatedAt:  k.CreatedAt,
		})
	}

	return res, total, nil
}

func (s *apiKeyService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *apiKeyService) Validate(ctx context.Context, rawKey string) (*entity.ApiKey, uint64, error) {
	hash := s.hashKey(rawKey)
	
	key, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid api key")
	}

	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return nil, 0, fmt.Errorf("api key expired")
	}

	// Calculate mask
	var mask uint64
	for _, p := range key.Role.Permissions {
		if p.ID <= 64 {
			mask |= (uint64(1) << (p.ID - 1))
		}
	}

	// Update last used asynchronously
	go s.repo.UpdateLastUsed(context.Background(), key.ID)

	return key, mask, nil
}

func (s *apiKeyService) hashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
