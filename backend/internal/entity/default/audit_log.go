package entity

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RequestID string    `gorm:"index" json:"request_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	UserEmail string    `json:"user_email"`
	ApiKeyID  *uint     `gorm:"index" json:"api_key_id"`
	Action    string    `gorm:"not null;index" json:"action"` // e.g. CREATE, UPDATE, DELETE
	Module    string    `gorm:"not null;index" json:"module"` // e.g. USER, ROLE, PRODUCT
	TargetID  string    `json:"target_id"`                    // ID of the resource affected
	Metadata  string    `gorm:"type:longtext" json:"metadata"` // JSON supporting information
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}
