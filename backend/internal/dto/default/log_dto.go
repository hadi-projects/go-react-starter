package dto

import "time"

type LogResponse struct {
	Level     string    `json:"level"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	UserID    *uint     `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Time      time.Time `json:"time"`
	Type      string    `json:"type"` // "auth" or "audit"
	RequestID string    `json:"request_id,omitempty"`
	Details   any       `json:"details,omitempty"`
}

type LogQuery struct {
	Type   string `form:"type"` // "auth", "audit", or "all"
	UserID uint   `form:"user_id"`
}
