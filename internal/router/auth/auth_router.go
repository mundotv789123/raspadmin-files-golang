package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/database/models"
)

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	LoginType string `json:"loginType" binding:"required"`
	Token     string `json:"token,omitempty"`
}

const refreshTokenExpireInMinutes = 7 * 24 * 60
const expireInSeconds = 15 * 60

func AuthLogin(c *gin.Context, db *gorm.DB) {
	if !config.AuthEnabled {
		c.JSON(400, gin.H{"message": "Authentication is disabled"})
		return
	}

	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(401, gin.H{"message": "Invalid request"})
		return
	}

	var session models.UserSession

	if strings.EqualFold(loginReq.LoginType, "CREDENTIALS") {
		if loginReq.Username == "" || loginReq.Password == "" {
			c.JSON(401, gin.H{"message": "Username and password are required"})
			return
		}
		if loginReq.Username != config.AppUsername || loginReq.Password != config.AppPassword {
			c.JSON(401, gin.H{"message": "Invalid username or password"})
			return
		}
		session = models.UserSession{}
		err := db.Create(&session).Error
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to create user session"})
			return
		}
	} else if strings.EqualFold(loginReq.LoginType, "REFRESH_TOKEN") {
		if loginReq.Token == "" {
			c.JSON(401, gin.H{"message": "Refresh token is required"})
			return
		}
		jwtDecoded, err := DecodeJwtToken(loginReq.Token)
		if err != nil {
			c.JSON(401, gin.H{"message": "Invalid token"})
			return
		}

		refreshToken, err := jwtDecoded.GetSubject()
		if err != nil {
			c.JSON(401, gin.H{"message": "Invalid token"})
			return
		}
		err = db.Where("refresh_token = ? AND expire_at > ?", refreshToken, int(time.Now().Unix())).First(&session).Error
		if err != nil {
			c.JSON(401, gin.H{"message": "Invalid refresh token"})
			return
		}
	} else {
		c.JSON(401, gin.H{"message": "Invalid login type"})
		return
	}

	session.SetRefreshToken(uuid.New().String(), refreshTokenExpireInMinutes)
	err := db.Save(session).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to save session"})
		return
	}

	expireAt := time.Now().Add(expireInSeconds * time.Second)
	claims := jwt.MapClaims{
		"username": loginReq.Username,
		"exp":      jwt.NewNumericDate(expireAt),
	}

	token, err := CreateJWTToken(claims)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create token"})
		return
	}

	refreshTokenExpireAt := time.Now().Add(refreshTokenExpireInMinutes * time.Minute)
	refreshTokenClaims := jwt.MapClaims{
		"sub": session.RefreshToken,
		"exp": jwt.NewNumericDate(refreshTokenExpireAt),
	}

	refreshToken, err := CreateJWTToken(refreshTokenClaims)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create refresh token"})
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
	})

	c.JSON(200, gin.H{
		"token":        token,
		"refreshToken": refreshToken,
	})
}
