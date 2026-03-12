package dto

import "time"

type WisudaResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWisudaRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateWisudaRequest struct {
	Name string `json:"name"`
}
