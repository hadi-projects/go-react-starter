package dto

import "time"

type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
