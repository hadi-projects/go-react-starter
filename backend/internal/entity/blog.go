package entity

import "time"

type Blog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
	Content string `gorm:"type:longtext" json:"content"`
	Thumbnail string `gorm:"type:varchar(255);not null" json:"thumbnail"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
