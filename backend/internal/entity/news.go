package entity

import "time"

type News struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
	Content string `gorm:"type:longtext" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
