package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/jwt"
)

func AuthenticationMiddleware(c *gin.Context) {
	if !config.AuthEnabled {
		c.Next()
		return
	}

	tokenString, _ := c.Cookie("token")

	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization token is missing"})
		return
	}

	_, err := jwt.DecodeJwtToken(tokenString)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	c.Next()
}
