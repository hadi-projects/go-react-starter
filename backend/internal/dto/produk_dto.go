package dto

import "time"

type ProdukResponse struct {
	ID        uint      `json:"id"`
	Name string `json:"name"`
	Harga int `json:"harga"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProdukRequest struct {
	Name string `json:"name" binding:"required"`
	Harga int `json:"harga" binding:"required"`
}

type UpdateProdukRequest struct {
	Name string `json:"name"`
	Harga int `json:"harga"`
}
