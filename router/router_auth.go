package router

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
	ijwt "github.com/mundotv789123/raspadmin/internal/jwt"
	"github.com/mundotv789123/raspadmin/internal/models"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Username and password are required")
	ErrInvalidLoginType          = errors.New("Invalid login type")
	ErrInvalidToken              = errors.New("Token is invalid")
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
		c.JSON(400, gin.H{"message": "Authentication disabled"})
		return
	}

	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	session, err := getSessionByRequest(loginReq, db)
	if err != nil {
		switch err {
		case ErrInvalidToken, ErrInvalidUsernameOrPassword:
			c.JSON(401, gin.H{"message": err.Error()})
		default:
			c.JSON(401, gin.H{"message": ErrInternalServerError.Error()})
		}
		return
	}

	session.SetRefreshToken(uuid.New().String(), refreshTokenExpireInMinutes)
	err = db.Save(session).Error
	if err != nil {
		c.JSON(500, gin.H{"message": ErrInternalServerError.Error()})
		return
	}

	expireAt := time.Now().Add(expireInSeconds * time.Second)
	claims := jwt.MapClaims{
		"username": loginReq.Username,
		"exp":      jwt.NewNumericDate(expireAt),
	}

	token, err := ijwt.CreateJWTToken(claims)
	if err != nil {
		c.JSON(500, gin.H{"message": ErrInternalServerError.Error()})
		return
	}

	refreshTokenExpireAt := time.Now().Add(refreshTokenExpireInMinutes * time.Minute)
	refreshTokenClaims := jwt.MapClaims{
		"sub": session.RefreshToken,
		"exp": jwt.NewNumericDate(refreshTokenExpireAt),
	}

	refreshToken, err := ijwt.CreateJWTToken(refreshTokenClaims)
	if err != nil {
		c.JSON(500, gin.H{"message": ErrInternalServerError.Error()})
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

func getSessionByRequest(loginReq LoginRequest, db *gorm.DB) (*models.UserSession, error) {
	if strings.EqualFold(loginReq.LoginType, "CREDENTIALS") {
		if loginReq.Username == "" || loginReq.Password == "" {
			return nil, ErrInvalidUsernameOrPassword
		}
		if loginReq.Username != config.AppUsername || loginReq.Password != config.AppPassword {
			return nil, ErrInvalidUsernameOrPassword
		}
		session := models.UserSession{}
		err := db.Create(&session).Error
		if err != nil {
			return nil, err
		}
		return &session, nil
	} else if strings.EqualFold(loginReq.LoginType, "REFRESH_TOKEN") {
		session, err := loginRefreshToken(loginReq, db)
		if err != nil {
			return nil, ErrInvalidToken
		}

		return session, nil
	}
	return nil, ErrInvalidLoginType
}

func loginRefreshToken(loginReq LoginRequest, db *gorm.DB) (*models.UserSession, error) {
	if loginReq.Token == "" {
		return nil, ErrInvalidToken
	}
	jwtDecoded, err := ijwt.DecodeJwtToken(loginReq.Token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	refreshToken, err := jwtDecoded.GetSubject()
	if err != nil {
		return nil, ErrInvalidToken
	}
	var session models.UserSession
	err = db.Where("refresh_token = ? AND expire_at > ?", refreshToken, int(time.Now().Unix())).First(&session).Error
	if err != nil || &session == nil {
		return nil, ErrInvalidToken
	}
	return &session, nil
}
