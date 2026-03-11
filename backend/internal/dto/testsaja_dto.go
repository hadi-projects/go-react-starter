package dto

import "time"

type TestsajaResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTestsajaRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTestsajaRequest struct {
	Name string `json:"name"`
}
