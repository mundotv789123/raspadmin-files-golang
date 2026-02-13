package router

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/router/middleware"
	"gorm.io/gorm"
)

var (
	ErrInternalServerError = errors.New("Internal server error")
)

type WebContext struct {
	DB *gorm.DB
}

func Index(c *gin.Context) {
	respose := gin.H{
		"message": "ping",
	}
	c.JSON(200, respose)
}

func (ctx *WebContext) Files(c *gin.Context) {
	filesResult, err := GetFiles(c.Query("path"), ctx.DB)
	if err != nil {
		if err == ErrFileNotFound {
			c.JSON(404, gin.H{"message": ErrFileNotFound})
			return
		}
		c.JSON(500, gin.H{"message": ErrInternalServerError})
		return
	}
	c.JSON(200, gin.H{"files": filesResult})
}

func (ctx *WebContext) AuthLogin(c *gin.Context) {
	AuthLogin(c, ctx.DB)
}

func (ctx *WebContext) Routers(r *gin.Engine) {
	r.Use(middleware.CorsMiddleware())

	apiRouter := r.Group("/api")

	apiRouter.GET("", Index)

	filesRouter := apiRouter.Group("/files")
	filesRouter.Use(middleware.AuthenticationMiddleware)
	filesRouter.GET("", ctx.Files)
	filesRouter.GET("open", OpenFile)

	authRouter := apiRouter.Group("/auth")
	authRouter.POST("/login", ctx.AuthLogin)
}
