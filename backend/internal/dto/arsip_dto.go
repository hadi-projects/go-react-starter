package dto

import "time"

type ArsipResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	Tanggal string `json:"tanggal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateArsipRequest struct {
	Name string `json:"name" binding:"required"`
	Tanggal string `json:"tanggal" binding:""`
}

type UpdateArsipRequest struct {
	Name string `json:"name"`
	Tanggal string `json:"tanggal"`
}
