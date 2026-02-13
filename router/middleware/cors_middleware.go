package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
			return slices.Contains(allowedOrigins, origin)
		}
		origin := c.Request.Header.Get("Origin")

		if isOriginAllowed(origin, config.AllowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
