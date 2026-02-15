package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/rs/zerolog"
)

var sensitiveKeys = map[string]bool{
	"password":              true,
	"token":                 true,
	"access_token":          true,
	"refresh_token":         true,
	"new_password":          true,
	"old_password":          true,
	"password_confirmation": true,
	"credit_card":           true,
	"cvv":                   true,
	"pan":                   true,
	"card_number":           true,
	"card_expiry":           true,
	"otp":                   true,
	"otp_code":              true,
	"secret_key":            true,
	"nik":                   true,
	"ktp_number":            true,
	"identity_number":       true,
}

var partialSensitiveKeys = map[string]bool{
	"email":        true,
	"phone":        true,
	"phone_number": true,
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()

		requestID := uuid.New().String()

		ctx.Set("request_id", requestID)
		ctx.Header("X-Request-ID", requestID)

		var body []byte
		if ctx.Request.Body != nil {
			body, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()

		latency := time.Since(start)
		statusCode := ctx.Writer.Status()
		userID, userExists := ctx.Get("user_id")

		var logEvent *zerolog.Event

		if strings.HasPrefix(path, "/api/auth") {
			if statusCode >= 500 {
				logEvent = logger.AuthLogger.Error()
			} else {
				logEvent = logger.AuthLogger.Info()
			}
		} else {
			if statusCode >= 500 {
				logEvent = logger.SystemLogger.Error()
			} else {
				logEvent = logger.SystemLogger.Info()
			}
		}

		logEvent.
			Str("request_id", requestID).
			Str("timestamp", time.Now().Format(time.RFC3339)).
			Uint("latency", uint(latency))
		if userExists {
			logEvent.Uint("user_id", userID.(uint))
		}
		requestDict := zerolog.Dict().
			Str("method", method).
			Str("path", path).
			Str("ip", clientIP).
			Str("user_agent", userAgent)

		if len(body) > 0 && json.Valid(body) {
			censoredBody := censorBody(body)
			requestDict.RawJSON("body", censoredBody)
		}

		responseDict := zerolog.Dict().
			Int("status_code", statusCode).
			Dur("latency", latency)

		if blw.body.Len() > 0 {
			resBody := blw.body.Bytes()
			if json.Valid(resBody) {
				censoredBody := censorBody(resBody)
				responseDict.RawJSON("body", censoredBody)
			}
		}

		logEvent.Dict("request", requestDict)
		logEvent.Dict("response", responseDict)
		logEvent.Msg("incoming request")

	}
}

func censorBody(body []byte) []byte {
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return body
	}

	maskedData := maskSensitiveData(data)
	if maskedBody, err := json.Marshal(maskedData); err == nil {
		return maskedBody
	}

	return body
}

func maskSensitiveData(data any) any {

	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			lowerKey := strings.ToLower(key)
			if sensitiveKeys[lowerKey] {
				v[key] = "***"
			} else if partialSensitiveKeys[lowerKey] {
				if strVal, ok := val.(string); ok {
					if strings.Contains(lowerKey, "email") {
						v[key] = maskEmail(strVal)
					} else {
						v[key] = maskPhone(strVal)
					}
				}
			} else {
				v[key] = maskSensitiveData(val)
			}
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = maskSensitiveData(val)
		}

		return v
	}
	return data
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local := parts[0]
	if len(local) <= 2 {
		return local[:1] + "***" + parts[1]
	}

	return "***@" + parts[1]
}

func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}

	return phone[:4] + "***" + phone[len(phone)-3:]
}
