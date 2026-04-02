package service

import (
	"context"
	"fmt"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type SettingService interface {
	GetSettings(ctx context.Context, category string) ([]dto.SettingResponse, error)
	GetConfigValue(ctx context.Context, key string) string
	BulkUpdate(ctx context.Context, req dto.UpdateSettingsRequest) error
}

type settingService struct {
	repo   repository.SettingRepository
	config *config.Config
}

func NewSettingService(repo repository.SettingRepository, cfg *config.Config) SettingService {
	return &settingService{repo: repo, config: cfg}
}

func (s *settingService) GetSettings(ctx context.Context, category string) ([]dto.SettingResponse, error) {
	var settings []entity.Setting
	var err error
	if category != "" {
		settings, err = s.repo.FindByCategory(ctx, category)
	} else {
		settings, err = s.repo.FindAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	var res []dto.SettingResponse
	for _, setting := range settings {
		res = append(res, dto.SettingResponse{
			Key:         setting.Key,
			Value:       setting.Value,
			Category:    setting.Category,
			FieldType:   setting.FieldType,
			Label:       setting.Label,
			Description: setting.Description,
		})
	}

	return res, nil
}

func (s *settingService) GetConfigValue(ctx context.Context, key string) string {
	setting, err := s.repo.FindByKey(ctx, key)
	if err == nil && setting != nil && setting.Value != "" {
		return setting.Value
	}

	// Fallbacks to .env config
	switch key {
	case "app_name":
		return s.config.App.Name
	case "smtp_host":
		return s.config.Mail.Host
	case "smtp_port":
		return fmt.Sprintf("%d", s.config.Mail.Port)
	case "smtp_user":
		return s.config.Mail.User
	case "smtp_pass":
		return s.config.Mail.Password
	case "smtp_from_name":
		return s.config.Mail.FromAddress // Assuming from address is the default name if not separate
	case "smtp_from_email":
		return s.config.Mail.FromAddress
	case "jwt_secret":
		return s.config.JWT.Secret
	case "jwt_issuer":
		return s.config.JWT.Issuer
	case "jwt_access_expiration":
		return s.config.JWT.AccessExpirationTime
	case "jwt_refresh_expiration":
		return s.config.JWT.RefreshExpirationTime
	case "cors_allowed_origins":
		return s.config.CORS.AllowedOrigins
	case "rate_limit_rps":
		return fmt.Sprintf("%d", s.config.RateLimit.Rps)
	case "rate_limit_burst":
		return fmt.Sprintf("%d", s.config.RateLimit.Burst)
	case "storage_base_path":
		return s.config.Storage.BasePath
	case "storage_max_file_size_mb":
		return fmt.Sprintf("%d", s.config.Storage.MaxFileSizeMB)
	case "db_host":
		return s.config.Database.Host
	case "db_port":
		return s.config.Database.Port
	case "db_username":
		return s.config.Database.UserName
	case "db_password":
		return s.config.Database.Password
	case "db_name":
		return s.config.Database.Name
	case "redis_host":
		return s.config.Redis.Host
	case "redis_port":
		return s.config.Redis.Port
	case "redis_password":
		return s.config.Redis.Password
	case "kafka_brokers":
		if len(s.config.Kafka.Brokers) > 0 {
			return s.config.Kafka.Brokers[0] // Simplified for single string return
		}
		return ""
	case "kafka_topic":
		return s.config.Kafka.Topic
	default:
		return ""
	}
}

func (s *settingService) BulkUpdate(ctx context.Context, req dto.UpdateSettingsRequest) error {
	var settingsToUpdate []entity.Setting
	for k, v := range req.Settings {
		settingsToUpdate = append(settingsToUpdate, entity.Setting{
			Key:   k,
			Value: v,
		})
	}

	if err := s.repo.BulkUpdate(ctx, settingsToUpdate); err != nil {
		return err
	}

	// Audit log
	userID, _ := ctx.Value(logger.CtxKeyUserID).(uint)
	logger.LogAudit(ctx, "UPDATE_SETTINGS", "SETTING", fmt.Sprintf("%d", userID), "Bulk update settings")

	return nil
}
