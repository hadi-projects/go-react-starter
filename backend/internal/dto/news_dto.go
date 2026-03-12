package dto

import "time"

type NewsResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNewsRequest struct {
	Name string `json:"name" binding:"required"`
	Content string `json:"content" binding:""`
}

type UpdateNewsRequest struct {
	Name string `json:"name"`
	Content string `json:"content"`
}
