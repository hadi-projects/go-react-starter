package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	sqlInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union.*select)|(insert.*into)|(delete.*from)|(drop.*table)|(update.*set)`),
		regexp.MustCompile(`(?i)(exec|execute)\s+`),
		regexp.MustCompile(`--|\#|\/\*|\*\/`),
		regexp.MustCompile(`(?i)(or|and)\s+[\d\w]+\s*=\s*[\d\w]+`),
	}
)

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSQLInjection(value) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Potentially malicious input detected"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

func containsSQLInjection(input string) bool {
	input = strings.ToLower(input)
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}
