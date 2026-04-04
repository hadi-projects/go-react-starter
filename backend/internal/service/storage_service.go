package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/storage"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	settingService "github.com/hadi-projects/go-react-starter/internal/service/default"
)

// ErrForbidden is returned when a share link access is denied.
var ErrForbidden = errors.New("forbidden")

// ErrNotFound is returned when an entity is not found.
var ErrNotFound = errors.New("not found")

// ErrInvalidPassword is returned when share link password is wrong.
var ErrInvalidPassword = errors.New("invalid password")

const (
	defaultMaxFileSizeMB    = 50
	cacheFileTTL            = 5 * time.Minute
	cacheShareLinkTTL       = 2 * time.Minute
)

type StorageService interface {
	// Authenticated operations
	Upload(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, description string) (*dto.StorageFileResponse, error)
	GetFiles(ctx context.Context, userID uint, pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error)
	GetFileByID(ctx context.Context, id, userID uint) (*dto.StorageFileResponse, error)
	DeleteFile(ctx context.Context, id, userID uint) error
	GetFileForDownload(ctx context.Context, id, userID uint) (io.ReadCloser, *entity.StorageFile, error)

	// Share link operations (authenticated)
	CreateShareLink(ctx context.Context, fileID, userID uint, req dto.CreateShareLinkRequest) (*dto.ShareLinkResponse, error)
	GetShareLinks(ctx context.Context, fileID, userID uint) ([]dto.ShareLinkResponse, error)
	UpdateShareLink(ctx context.Context, shareLinkID, userID uint, req dto.UpdateShareLinkRequest) (*dto.ShareLinkResponse, error)
	RevokeShareLink(ctx context.Context, shareLinkID, userID uint) error
	GetShareLinkLogs(ctx context.Context, shareLinkID, userID uint) ([]dto.ShareLinkAccessResponse, error)

	// Public access (no auth)
	GetPublicFileInfo(ctx context.Context, token string) (*dto.PublicFileResponse, error)
	ServePublicFile(ctx context.Context, token, password, ip, userAgent string) (io.ReadCloser, *entity.StorageFile, bool, error)
	ServePublicSystemFile(ctx context.Context, id uint) (io.ReadCloser, *entity.StorageFile, error)
}

type storageService struct {
	fileRepo      repository.StorageFileRepository
	shareLinkRepo repository.ShareLinkRepository
	driver        storage.Driver
	cache          cache.CacheService
	frontendURL    string
	settingService settingService.SettingService
}

func NewStorageService(
	fileRepo repository.StorageFileRepository,
	shareLinkRepo repository.ShareLinkRepository,
	driver storage.Driver,
	cacheService cache.CacheService,
	frontendURL string,
	settingSvc settingService.SettingService,
) StorageService {
	return &storageService{
		fileRepo:       fileRepo,
		shareLinkRepo:  shareLinkRepo,
		driver:         driver,
		cache:          cacheService,
		frontendURL:    frontendURL,
		settingService: settingSvc,
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func detectMIME(f multipart.File) (string, error) {
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	mime := http.DetectContentType(buf[:n])
	// Reset read position
	if seeker, ok := f.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}
	return mime, nil
}

func (s *storageService) shareURL(token string) string {
	return strings.TrimRight(s.frontendURL, "/") + "/share/" + token
}

func (s *storageService) mapFileToResponse(ctx context.Context, f *entity.StorageFile) *dto.StorageFileResponse {
	shareCount, _ := s.fileRepo.CountShareLinks(ctx, f.ID)
	return &dto.StorageFileResponse{
		ID:           f.ID,
		UserID:       f.UserID,
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		Size:         f.Size,
		SizeHuman:    formatSize(f.Size),
		Description:  f.Description,
		ShareCount:   shareCount,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
}

func (s *storageService) mapLinkToResponse(link *entity.ShareLink) *dto.ShareLinkResponse {
	return &dto.ShareLinkResponse{
		ID:            link.ID,
		FileID:        link.FileID,
		Token:         link.Token,
		Label:         link.Label,
		ShareURL:      s.shareURL(link.Token),
		AccessType:    string(link.AccessType),
		MaxViews:      link.MaxViews,
		ViewCount:     link.ViewCount,
		ExpiresAt:     link.ExpiresAt,
		HasPassword:   link.PasswordHash != nil,
		AllowDownload: link.AllowDownload,
		IsActive:      link.IsActive,
		CreatedAt:     link.CreatedAt,
		UpdatedAt:     link.UpdatedAt,
	}
}

// validateLink checks all access conditions. Returns ErrForbidden or ErrInvalidPassword on failure.
func (s *storageService) validateLink(link *entity.ShareLink, password string) error {
	if !link.IsActive {
		return ErrForbidden
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return ErrForbidden
	}
	if link.AccessType == entity.AccessTypeLimited && link.MaxViews != nil {
		if link.ViewCount >= *link.MaxViews {
			return ErrForbidden
		}
	}
	if link.PasswordHash != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(*link.PasswordHash), []byte(password)); err != nil {
			return ErrInvalidPassword
		}
	}
	return nil
}

// consumeLink increments view count and deactivates one_time links.
func (s *storageService) consumeLink(ctx context.Context, link *entity.ShareLink, ip, userAgent string) {
	link.ViewCount++
	if link.AccessType == entity.AccessTypeOneTime {
		link.IsActive = false
	}
	s.shareLinkRepo.Update(ctx, link)
	s.shareLinkRepo.RecordAccess(ctx, &entity.ShareLinkAccess{
		ShareLinkID: link.ID,
		IPAddress:   ip,
		UserAgent:   userAgent,
	})
	s.cache.Delete(ctx, fmt.Sprintf("share_link:token:%s", link.Token))
}

// ─── File operations ──────────────────────────────────────────────────────────

func (s *storageService) Upload(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, description string) (*dto.StorageFileResponse, error) {
	// Size validation
	maxSizeValue := s.settingService.GetConfigValue(ctx, "storage_max_file_size_mb")
	maxMB, _ := strconv.ParseInt(maxSizeValue, 10, 64)
	if maxMB <= 0 {
		maxMB = defaultMaxFileSizeMB
	}
	maxBytes := maxMB * 1024 * 1024
	if fileHeader.Size > maxBytes {
		return nil, fmt.Errorf("file size %s exceeds maximum allowed %s",
			formatSize(fileHeader.Size), formatSize(maxBytes))
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Detect MIME from content (not extension) for security
	mimeType, err := detectMIME(src)
	if err != nil {
		return nil, fmt.Errorf("failed to detect file type: %w", err)
	}

	// Build stored name and path
	ext := ""
	if idx := strings.LastIndex(fileHeader.Filename, "."); idx >= 0 {
		ext = fileHeader.Filename[idx:]
	}
	storedName := uuid.New().String() + ext
	now := time.Now()
	storagePath := fmt.Sprintf("%d/%02d/", now.Year(), now.Month())
	key := storagePath + storedName

	// Save to storage driver
	if err := s.driver.Save(ctx, key, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Persist metadata
	fileEntity := &entity.StorageFile{
		UserID:       userID,
		OriginalName: fileHeader.Filename,
		StoredName:   storedName,
		StoragePath:  storagePath,
		MimeType:     mimeType,
		Size:         fileHeader.Size,
		Description:  description,
	}
	if err := s.fileRepo.Create(ctx, fileEntity); err != nil {
		// Best-effort cleanup
		s.driver.Delete(ctx, key)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	s.cache.DeletePattern(ctx, fmt.Sprintf("storage:user:%d:*", userID))
	logger.LogAudit(ctx, "UPLOAD", "STORAGE_FILE", fmt.Sprintf("%d", fileEntity.ID),
		fmt.Sprintf("name: %s, size: %s, mime: %s", fileEntity.OriginalName, formatSize(fileEntity.Size), mimeType))

	return s.mapFileToResponse(ctx, fileEntity), nil
}

func (s *storageService) GetFiles(ctx context.Context, userID uint, pagination *defaultDto.PaginationRequest) (*defaultDto.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("storage:user:%d:page:%d:limit:%d:search:%s",
		userID, pagination.GetPage(), pagination.GetLimit(), pagination.Search)

	var cached defaultDto.PaginationResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	files, total, err := s.fileRepo.FindAll(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	var responses []dto.StorageFileResponse
	for _, f := range files {
		responses = append(responses, *s.mapFileToResponse(ctx, &f))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.GetLimit())))
	result := &defaultDto.PaginationResponse{
		Data: responses,
		Meta: defaultDto.PaginationMeta{
			CurrentPage: pagination.GetPage(),
			TotalPages:  totalPages,
			TotalData:   total,
			Limit:       pagination.GetLimit(),
		},
	}

	s.cache.Set(ctx, cacheKey, result, cacheFileTTL)
	return result, nil
}

func (s *storageService) GetFileByID(ctx context.Context, id, userID uint) (*dto.StorageFileResponse, error) {
	file, err := s.fileRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.mapFileToResponse(ctx, file), nil
}

func (s *storageService) DeleteFile(ctx context.Context, id, userID uint) error {
	file, err := s.fileRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return ErrNotFound
	}

	key := file.StoragePath + file.StoredName
	if err := s.driver.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	if err := s.fileRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.cache.DeletePattern(ctx, fmt.Sprintf("storage:user:%d:*", userID))
	logger.LogAudit(ctx, "DELETE", "STORAGE_FILE", fmt.Sprintf("%d", id),
		fmt.Sprintf("name: %s", file.OriginalName))
	return nil
}

func (s *storageService) GetFileForDownload(ctx context.Context, id, userID uint) (io.ReadCloser, *entity.StorageFile, error) {
	file, err := s.fileRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, nil, ErrNotFound
	}
	reader, err := s.driver.Get(ctx, file.StoragePath+file.StoredName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	return reader, file, nil
}

// ─── Share link operations ────────────────────────────────────────────────────

func (s *storageService) CreateShareLink(ctx context.Context, fileID, userID uint, req dto.CreateShareLinkRequest) (*dto.ShareLinkResponse, error) {
	// Verify ownership
	if _, err := s.fileRepo.FindByIDAndUserID(ctx, fileID, userID); err != nil {
		return nil, ErrNotFound
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	link := &entity.ShareLink{
		FileID:        fileID,
		Token:         token,
		Label:         req.Label,
		AccessType:    entity.AccessType(req.AccessType),
		MaxViews:      req.MaxViews,
		ExpiresAt:     req.ExpiresAt,
		AllowDownload: req.AllowDownload,
		IsActive:      true,
	}

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		h := string(hash)
		link.PasswordHash = &h
	}

	if err := s.shareLinkRepo.Create(ctx, link); err != nil {
		return nil, err
	}

	s.cache.DeletePattern(ctx, fmt.Sprintf("storage:user:%d:*", userID))
	logger.LogAudit(ctx, "CREATE", "SHARE_LINK", fmt.Sprintf("%d", link.ID),
		fmt.Sprintf("file_id: %d, access_type: %s", fileID, req.AccessType))

	return s.mapLinkToResponse(link), nil
}

func (s *storageService) GetShareLinks(ctx context.Context, fileID, userID uint) ([]dto.ShareLinkResponse, error) {
	if _, err := s.fileRepo.FindByIDAndUserID(ctx, fileID, userID); err != nil {
		return nil, ErrNotFound
	}
	links, err := s.shareLinkRepo.FindByFileID(ctx, fileID)
	if err != nil {
		return nil, err
	}
	var responses []dto.ShareLinkResponse
	for _, l := range links {
		responses = append(responses, *s.mapLinkToResponse(&l))
	}
	return responses, nil
}

func (s *storageService) UpdateShareLink(ctx context.Context, shareLinkID, userID uint, req dto.UpdateShareLinkRequest) (*dto.ShareLinkResponse, error) {
	link, err := s.shareLinkRepo.FindByIDAndUserID(ctx, shareLinkID, userID)
	if err != nil {
		return nil, ErrNotFound
	}

	if req.Label != nil {
		link.Label = *req.Label
	}
	if req.AccessType != nil {
		link.AccessType = entity.AccessType(*req.AccessType)
	}
	if req.MaxViews != nil {
		link.MaxViews = req.MaxViews
	}
	if req.ExpiresAt != nil {
		link.ExpiresAt = req.ExpiresAt
	}
	if req.AllowDownload != nil {
		link.AllowDownload = *req.AllowDownload
	}
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}
	if req.Password != nil {
		if *req.Password == "" {
			link.PasswordHash = nil
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			h := string(hash)
			link.PasswordHash = &h
		}
	}

	if err := s.shareLinkRepo.Update(ctx, link); err != nil {
		return nil, err
	}

	s.cache.Delete(ctx, fmt.Sprintf("share_link:token:%s", link.Token))
	logger.LogAudit(ctx, "UPDATE", "SHARE_LINK", fmt.Sprintf("%d", link.ID), "")
	return s.mapLinkToResponse(link), nil
}

func (s *storageService) RevokeShareLink(ctx context.Context, shareLinkID, userID uint) error {
	link, err := s.shareLinkRepo.FindByIDAndUserID(ctx, shareLinkID, userID)
	if err != nil {
		return ErrNotFound
	}
	if err := s.shareLinkRepo.Delete(ctx, shareLinkID); err != nil {
		return err
	}
	s.cache.Delete(ctx, fmt.Sprintf("share_link:token:%s", link.Token))
	logger.LogAudit(ctx, "DELETE", "SHARE_LINK", fmt.Sprintf("%d", shareLinkID), "")
	return nil
}

func (s *storageService) GetShareLinkLogs(ctx context.Context, shareLinkID, userID uint) ([]dto.ShareLinkAccessResponse, error) {
	if _, err := s.shareLinkRepo.FindByIDAndUserID(ctx, shareLinkID, userID); err != nil {
		return nil, ErrNotFound
	}
	logs, err := s.shareLinkRepo.GetAccessLogs(ctx, shareLinkID)
	if err != nil {
		return nil, err
	}
	var responses []dto.ShareLinkAccessResponse
	for _, l := range logs {
		responses = append(responses, dto.ShareLinkAccessResponse{
			ID:          l.ID,
			ShareLinkID: l.ShareLinkID,
			IPAddress:   l.IPAddress,
			UserAgent:   l.UserAgent,
			AccessedAt:  l.AccessedAt,
		})
	}
	return responses, nil
}

// ─── Public access ────────────────────────────────────────────────────────────

func (s *storageService) GetPublicFileInfo(ctx context.Context, token string) (*dto.PublicFileResponse, error) {
	link, err := s.shareLinkRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, ErrNotFound
	}

	// Check basic validity (without consuming)
	if !link.IsActive {
		return nil, ErrForbidden
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, ErrForbidden
	}
	if link.AccessType == entity.AccessTypeLimited && link.MaxViews != nil {
		if link.ViewCount >= *link.MaxViews {
			return nil, ErrForbidden
		}
	}

	return &dto.PublicFileResponse{
		Token:            link.Token,
		Label:            link.Label,
		OriginalName:     link.File.OriginalName,
		MimeType:         link.File.MimeType,
		Size:             link.File.Size,
		SizeHuman:        formatSize(link.File.Size),
		AccessType:       string(link.AccessType),
		ViewCount:        link.ViewCount,
		MaxViews:         link.MaxViews,
		ExpiresAt:        link.ExpiresAt,
		AllowDownload:    link.AllowDownload,
		RequiresPassword: link.PasswordHash != nil,
	}, nil
}

func (s *storageService) ServePublicFile(ctx context.Context, token, password, ip, userAgent string) (io.ReadCloser, *entity.StorageFile, bool, error) {
	link, err := s.shareLinkRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, nil, false, ErrNotFound
	}

	if err := s.validateLink(link, password); err != nil {
		return nil, nil, false, err
	}

	reader, err := s.driver.Get(ctx, link.File.StoragePath+link.File.StoredName)
	if err != nil {
		return nil, nil, false, fmt.Errorf("file not available: %w", err)
	}

	// Consume the link (async to not block response)
	go s.consumeLink(context.Background(), link, ip, userAgent)

	return reader, &link.File, link.AllowDownload, nil
}

func (s *storageService) ServePublicSystemFile(ctx context.Context, id uint) (io.ReadCloser, *entity.StorageFile, error) {
	// Check if ID is logo or favicon in settings
	logoID := s.settingService.GetConfigValue(ctx, "app_logo")
	faviconID := s.settingService.GetConfigValue(ctx, "app_favicon")

	idStr := strconv.FormatUint(uint64(id), 10)
	if idStr != logoID && idStr != faviconID {
		return nil, nil, ErrForbidden
	}

	file, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, ErrNotFound
	}

	reader, err := s.driver.Get(ctx, file.StoragePath+file.StoredName)
	if err != nil {
		return nil, nil, fmt.Errorf("file not available: %w", err)
	}

	return reader, file, nil
}
