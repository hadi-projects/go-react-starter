package dto

import "time"

type BlogResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	Content string `json:"content"`
	Thumbnail string `json:"thumbnail"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBlogRequest struct {
	Name string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
	Thumbnail string `json:"thumbnail" binding:""`
}

type UpdateBlogRequest struct {
	Name string `json:"name"`
	Content string `json:"content"`
	Thumbnail string `json:"thumbnail"`
}
