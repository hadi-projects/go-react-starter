package dto

import "time"

type AdminResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAdminRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:""`
}

type UpdateAdminRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
}
