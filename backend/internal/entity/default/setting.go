package entity

import "time"

type Setting struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value       string    `gorm:"type:longtext" json:"value"`
	Category    string    `gorm:"type:varchar(50);index;not null" json:"category"`
	FieldType   string    `gorm:"type:varchar(20);not null;default:'text'" json:"field_type"` // text, number, boolean, file
	Label       string    `gorm:"type:varchar(255)" json:"label"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
