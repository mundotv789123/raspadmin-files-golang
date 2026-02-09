package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/router/files"
)

func Index(c *gin.Context) {
	respose := gin.H{
		"message": "Hello, World!",
	}
	c.JSON(200, respose)
}

func Files(c *gin.Context) {
	files, err := files.GetFiles(c.Query("path"))
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"files": files})
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
