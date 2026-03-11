package dto

import "time"

type CookResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCookRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCookRequest struct {
	Name string `json:"name"`
}
