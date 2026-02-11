package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/router/files"
	"gorm.io/gorm"
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
	filesResult, err := files.GetFiles(c.Query("path"), ctx.DB)
	if err != nil {
		if err == files.ErrFileNotFound {
			c.JSON(404, gin.H{"message": "File or directory not found"})
			return
		}
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"files": filesResult})
}

func OpenFile(c *gin.Context) {
	file, err := files.OpenFile(c.Query("path"))
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	c.File(file.Name())
}

func (ctx *WebContext) Routers(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.GET("", Index)
	apiRouter.GET("files", ctx.Files)
	apiRouter.GET("files/open", OpenFile)
}
