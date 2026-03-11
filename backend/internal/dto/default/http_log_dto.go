package dto

import "time"

type HttpLogQuery struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	Method     string `form:"method"`
	Path       string `form:"path"`
	StatusCode int    `form:"status_code"`
}

func (q *HttpLogQuery) GetOffset() int {
	return (q.GetPage() - 1) * q.GetLimit()
}

func (q *HttpLogQuery) GetLimit() int {
	if q.Limit == 0 {
		return 10
	}
	return q.Limit
}

func (q *HttpLogQuery) GetPage() int {
	if q.Page == 0 {
		return 1
	}
	return q.Page
}

type HttpLogResponse struct {
	ID              uint      `json:"id"`
	RequestID       string    `json:"request_id"`
	Method          string    `json:"method"`
	Path            string    `json:"path"`
	ClientIP        string    `json:"client_ip"`
	UserAgent       string    `json:"user_agent"`
	RequestHeaders  string    `json:"request_headers"`
	RequestBody     string    `json:"request_body"`
	StatusCode      int       `json:"status_code"`
	ResponseHeaders string    `json:"response_headers"`
	ResponseBody    string    `json:"response_body"`
	Latency         int64     `json:"latency"`
	UserID          *uint     `json:"user_id,omitempty"`
	UserEmail       string    `json:"user_email,omitempty"`
	MiddlewareTrace string    `json:"middleware_trace,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
