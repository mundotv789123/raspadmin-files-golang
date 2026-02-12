package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
)

type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	LoginType string `json:"loginType" binding:"required"`
	Token     string `json:"token,omitempty"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

const expireInSeconds = 24 * 60 * 60

func AuthLogin(c *gin.Context) {
	if !config.AuthEnabled {
		c.JSON(400, gin.H{"message": "Authentication is disabled"})
		return
	}

	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(401, gin.H{"message": "Invalid request"})
		return
	}

	if loginReq.Username != config.AppUsername || loginReq.Password != config.AppPassword {
		c.JSON(401, gin.H{"message": "Invalid username or password"})
		return
	}

	expireAt := time.Now().Add(expireInSeconds * time.Second)
	claims := jwt.MapClaims{
		"username": loginReq.Username,
		"exp":      jwt.NewNumericDate(expireAt), //TODO implementar refresh token
	}

	token, erro := createJWTToken(claims)
	if erro != nil {
		c.JSON(500, gin.H{"message": "Failed to create token"})
		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   "",
		Expires:  expireAt,
		MaxAge:   expireInSeconds,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Partitioned: true, // Go 1.22+
	})

	c.JSON(200, gin.H{
		"token": token,
	})
}

func createJWTToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
