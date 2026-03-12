package dto

import "time"

type MainnnResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	Makananan string `json:"makananan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateMainnnRequest struct {
	Name string `json:"name" binding:"required"`
	Makananan string `json:"makananan" binding:""`
}

type UpdateMainnnRequest struct {
	Name string `json:"name"`
	Makananan string `json:"makananan"`
}
