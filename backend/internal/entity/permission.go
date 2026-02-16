package entity

import "time"

type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"unique;type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500);default:'default'" json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
