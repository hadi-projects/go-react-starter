package entity

import "time"

type ShareLinkAccess struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShareLinkID uint      `gorm:"not null;index" json:"share_link_id"`
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string    `gorm:"type:varchar(500)" json:"user_agent"`
	AccessedAt  time.Time `gorm:"autoCreateTime" json:"accessed_at"`
}
