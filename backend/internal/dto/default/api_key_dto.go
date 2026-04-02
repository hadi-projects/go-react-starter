package dto

import "time"

type ApiKeyCreateRequest struct {
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=uuid sk_tp"`
	RoleID        uint   `json:"role_id" binding:"required"`
	ExpiresInDays int    `json:"expires_in_days"`
	AllowedIPs    string `json:"allowed_ips"`
}

type ApiKeyCreateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	RawKey    string    `json:"raw_key"`
	Prefix    string    `json:"prefix"`
	RoleName  string    `json:"role_name"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type ApiKeyResponse struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	RoleName   string     `json:"role_name"`
	AllowedIPs string     `json:"allowed_ips"`
	ExpiresAt  *time.Time `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
