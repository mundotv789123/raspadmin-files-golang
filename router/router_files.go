package router

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/repository"
	"github.com/mundotv789123/raspadmin/router/dto"
	"github.com/mundotv789123/raspadmin/router/system"
	"gorm.io/gorm"
)

var hiddenFilesRegex = regexp.MustCompile(`^[\._].*$`)

var (
	ErrFileNotFound = errors.New("File or directory not found")
)

func GetFiles(path string, db *gorm.DB) ([]dto.FileDto, error) {
	path, err := system.SafeJoinPath(path)
	if err != nil {
		return nil, ErrFileNotFound
	}

	fileState, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	if !fileState.IsDir() {
		filePath := path[len(config.AbsRootDir):]
		return []dto.FileDto{dto.NewFileDto(fileState, filePath, true, "", path)}, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	filesList := make([]dto.FileDto, 0)
	parentPath := path[len(config.AbsRootDir):]

	if parentPath == "" {
		parentPath = "/"
	}

	filesDb, err := repository.GetFilesMapFromParentPath(db, parentPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if hiddenFilesRegex.MatchString(file.Name()) {
			continue
		}

		fileState, err := file.Info()
		if err != nil {
			continue
		}

		fileDb, ok := filesDb[file.Name()]
		var fileIcon string
		if ok && fileDb.IconPath != nil {
			fileIcon = *fileDb.IconPath
		}

		filePath := filepath.Join(parentPath, file.Name())
		fileDto := dto.NewFileDto(
			fileState,
			filePath,
			false,
			fileIcon,
			filepath.Join(config.AbsRootDir, filePath),
		)
		filesList = append(filesList, fileDto)
	}

	return filesList, nil
}

func OpenFile(c *gin.Context) {
	filePath, err := system.SafeJoinPath(c.Query("path"))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"message": ErrFileNotFound.Error()})
			return
		}
		c.JSON(500, gin.H{"message": ErrInternalServerError.Error()})
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	if file == nil {
		c.JSON(404, gin.H{"message": ErrFileNotFound.Error()})
		return
	}
	defer file.Close()

	c.File(file.Name())
}
