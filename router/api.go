package router

import "github.com/gin-gonic/gin"

func Index(c *gin.Context) {
	respose := gin.H{
		"message": "Hello, World!",
	}
	c.JSON(200, respose)
}
