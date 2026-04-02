package entity

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"unique;type:varchar(50);not null" json:"name"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	Category    string       `gorm:"type:varchar(20);default:'user';not null" json:"category"` // user or api
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permissions []Permission `gorm:"many2many:role_has_permissions;" json:"permissions"`
}
