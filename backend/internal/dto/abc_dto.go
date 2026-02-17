package dto

import "time"

type AbcResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAbcRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateAbcRequest struct {
	Name string `json:"name"`
}
