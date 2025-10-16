package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setSafeHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Cache-Control", "no-cache;no-store;max-age=0;must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
}

func SafeHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		setSafeHeaders(c.Writer)
		c.Next()
	}
}
