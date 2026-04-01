package entity

import "time"

type AccessType string

const (
	AccessTypeOneTime   AccessType = "one_time"
	AccessTypeUnlimited AccessType = "unlimited"
	AccessTypeLimited   AccessType = "limited"
	AccessTypeTimed     AccessType = "timed"
)

type ShareLink struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID       uint        `gorm:"not null;index" json:"file_id"`
	File         StorageFile `gorm:"foreignKey:FileID" json:"-"`
	Token        string      `gorm:"type:varchar(64);not null;uniqueIndex" json:"token"`
	Label        string      `gorm:"type:varchar(100);default:''" json:"label"`
	AccessType   AccessType  `gorm:"type:enum('one_time','unlimited','limited','timed');default:'unlimited'" json:"access_type"`
	MaxViews     *int        `gorm:"default:null" json:"max_views"`
	ViewCount    int         `gorm:"default:0" json:"view_count"`
	ExpiresAt    *time.Time  `gorm:"default:null" json:"expires_at"`
	PasswordHash *string     `gorm:"type:varchar(255);default:null" json:"-"`
	HasPassword  bool        `gorm:"-" json:"has_password"`
	AllowDownload bool       `gorm:"default:true" json:"allow_download"`
	IsActive     bool        `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
