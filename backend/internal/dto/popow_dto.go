package dto

import "time"

type PopowResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePopowRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdatePopowRequest struct {
	Name string `json:"name"`
}
