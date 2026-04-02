package entity

import "time"

type ApiKey struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"` // Creator/Owner
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	KeyHash   string    `gorm:"unique;type:varchar(255);not null" json:"-"`
	Prefix    string    `gorm:"type:varchar(20);not null" json:"prefix"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	AllowedIPs string    `gorm:"type:text" json:"allowed_ips"` // Comma-separated
	ExpiresAt *time.Time `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
