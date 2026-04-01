package dto

import "time"

// ---- Request DTOs ----

// UploadFileRequest is bound from multipart/form-data.
// The actual file is extracted via c.FormFile("file") in the handler.
type UploadFileRequest struct {
	Description string `form:"description"`
}

type CreateShareLinkRequest struct {
	Label         string     `json:"label"`
	AccessType    string     `json:"access_type" binding:"required,oneof=one_time unlimited limited timed"`
	MaxViews      *int       `json:"max_views"`
	ExpiresAt     *time.Time `json:"expires_at"`
	Password      string     `json:"password"`
	AllowDownload bool       `json:"allow_download"`
}

type UpdateShareLinkRequest struct {
	Label         *string    `json:"label"`
	AccessType    *string    `json:"access_type" binding:"omitempty,oneof=one_time unlimited limited timed"`
	MaxViews      *int       `json:"max_views"`
	ExpiresAt     *time.Time `json:"expires_at"`
	Password      *string    `json:"password"`
	AllowDownload *bool      `json:"allow_download"`
	IsActive      *bool      `json:"is_active"`
}

// ---- Response DTOs ----

type StorageFileResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	SizeHuman    string    `json:"size_human"`
	Description  string    `json:"description"`
	ShareCount   int64     `json:"share_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ShareLinkResponse struct {
	ID            uint       `json:"id"`
	FileID        uint       `json:"file_id"`
	Token         string     `json:"token"`
	Label         string     `json:"label"`
	ShareURL      string     `json:"share_url"`
	AccessType    string     `json:"access_type"`
	MaxViews      *int       `json:"max_views"`
	ViewCount     int        `json:"view_count"`
	ExpiresAt     *time.Time `json:"expires_at"`
	HasPassword   bool       `json:"has_password"`
	AllowDownload bool       `json:"allow_download"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ShareLinkAccessResponse struct {
	ID          uint      `json:"id"`
	ShareLinkID uint      `json:"share_link_id"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	AccessedAt  time.Time `json:"accessed_at"`
}

// PublicFileResponse is returned for unauthenticated share link access.
// It deliberately excludes internal paths and sensitive data.
type PublicFileResponse struct {
	Token            string     `json:"token"`
	Label            string     `json:"label"`
	OriginalName     string     `json:"original_name"`
	MimeType         string     `json:"mime_type"`
	Size             int64      `json:"size"`
	SizeHuman        string     `json:"size_human"`
	AccessType       string     `json:"access_type"`
	ViewCount        int        `json:"view_count"`
	MaxViews         *int       `json:"max_views"`
	ExpiresAt        *time.Time `json:"expires_at"`
	AllowDownload    bool       `json:"allow_download"`
	RequiresPassword bool       `json:"requires_password"`
}
