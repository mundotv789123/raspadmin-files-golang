package router

import (
	"os"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

func Index(c *gin.Context) {
	respose := gin.H{
		"message": "Hello, World!",
	}
	c.JSON(200, respose)
}

func Files(c *gin.Context) {
	files, err := os.ReadDir("./files")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"message": "Files directory does not exist"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	filesList := make([]FileInfo, len(files))
	for i, file := range files {
		filesList[i] = FileInfo{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		}
	}
	c.JSON(200, gin.H{"files": filesList})
}
