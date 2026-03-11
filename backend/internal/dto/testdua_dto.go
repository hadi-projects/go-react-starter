package dto

import "time"

type TestduaResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTestduaRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTestduaRequest struct {
	Name string `json:"name"`
}
