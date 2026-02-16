package router

import (
	"errors"
	"os"
	"path/filepath"

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

const PUBLIC_DIR = "./public"

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
	loadPublicDir(r)

	apiRouter := r.Group("/api")

	apiRouter.GET("", Index)

	filesRouter := apiRouter.Group("/files")
	filesRouter.Use(middleware.AuthenticationMiddleware)
	filesRouter.GET("", ctx.Files)
	filesRouter.GET("open", OpenFile)

	authRouter := apiRouter.Group("/auth")
	authRouter.POST("/login", ctx.AuthLogin)
}

func loadPublicDir(r *gin.Engine) {
	if files, err := os.ReadDir(PUBLIC_DIR); err == nil {
		for _, file := range files {
			fileName := file.Name()
			if fileName == "index.html" {
				r.GET("", func(c *gin.Context) {
					c.File(filepath.Join(PUBLIC_DIR, fileName))
				})
				continue
			}
			r.Static(fileName, filepath.Join(PUBLIC_DIR, fileName))
		}
	}
}
