package entity

import "time"

type StorageFile struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	StoredName   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"stored_name"`
	StoragePath  string    `gorm:"type:varchar(500);not null" json:"storage_path"`
	MimeType     string    `gorm:"type:varchar(100);not null" json:"mime_type"`
	Size         int64     `gorm:"not null" json:"size"`
	Description  string    `gorm:"type:varchar(500);default:''" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
