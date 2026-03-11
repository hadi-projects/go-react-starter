package entity

import "time"

type HttpLog struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RequestID       string    `gorm:"not null" json:"request_id"`
	Method          string    `gorm:"not null;index" json:"method"`
	Path            string    `gorm:"not null;type:text" json:"path"`
	ClientIP        string    `gorm:"not null" json:"client_ip"`
	UserAgent       string    `gorm:"type:text" json:"user_agent"`
	RequestHeaders  string    `gorm:"type:text" json:"request_headers"`
	RequestBody     string    `gorm:"type:longtext" json:"request_body"`
	StatusCode      int       `gorm:"not null;index" json:"status_code"`
	ResponseHeaders string    `gorm:"type:text" json:"response_headers"`
	ResponseBody    string    `gorm:"type:longtext" json:"response_body"`
	Latency         int64     `gorm:"not null" json:"latency"` // in milliseconds
	UserID          *uint     `gorm:"index" json:"user_id"`
	UserEmail       string    `json:"user_email"`
	MiddlewareTrace string    `gorm:"type:text" json:"middleware_trace"`
	CreatedAt       time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}
