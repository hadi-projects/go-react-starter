package dto

import "time"

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Category    string               `json:"category"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	Category      string `json:"category" binding:"oneof=user api"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Category      string `json:"category" binding:"oneof=user api"`
	PermissionIDs []uint `json:"permission_ids"`
}
