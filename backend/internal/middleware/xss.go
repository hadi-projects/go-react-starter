package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`<script[\s\S]*?>[\s\S]*?</script>`),
	regexp.MustCompile(`javascript:`),
	regexp.MustCompile(`on\w+\s*=`),
}

func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")

		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				if containsXSS(value) {
					values[i] = sanitizeXSS(value)
				}
			}
			c.Request.URL.Query()[key] = values
		}

		c.Next()
	}
}

func containsXSS(input string) bool {
	for _, pattern := range xssPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

func sanitizeXSS(input string) string {
	for _, pattern := range xssPatterns {
		input = pattern.ReplaceAllString(input, "")
	}
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")

	return input
}
