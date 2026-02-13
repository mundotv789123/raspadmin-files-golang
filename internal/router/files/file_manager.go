package files

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/database/models"
	"gorm.io/gorm"

	"github.com/djherbis/times"
)

var hiddenFilesRegex = regexp.MustCompile("^[\\._].*$")

type FileDto struct {
	Name      string    `json:"name"`
	IsDir     bool      `json:"is_dir"`
	Path      string    `json:"path"`
	Type      string    `json:"type,omitempty"`
	Icon      string    `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Open      bool      `json:"open,omitempty"`
}

func NewFileDto(file os.FileInfo, path string, open bool, iconPath string, fullPath string) FileDto {
	contentType := mime.TypeByExtension(filepath.Ext(file.Name()))

	var createdAt time.Time
	var updatedAt time.Time
	t, err := times.Stat(fullPath)
	if err == nil {
		createdAt = t.BirthTime()
		updatedAt = t.ModTime()
	}

	return FileDto{
		Name:      file.Name(),
		IsDir:     file.IsDir(),
		Path:      "/" + path,
		Type:      contentType,
		Icon:      iconPath,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Open:      open,
	}
}

var (
	ErrFileNotFound = errors.New("File or directory not found")
)

func GetFiles(path string, db *gorm.DB) ([]FileDto, error) {
	path, err := safeJoin(path)
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
		return []FileDto{NewFileDto(fileState, filePath, true, "", path)}, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	filesList := make([]FileDto, 0)
	parentPath := path[len(config.AbsRootDir):]

	filesDb, err := getDbFiles(db, parentPath)
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
		if ok {
			fileIcon = *fileDb.IconPath
		}

		filePath := filepath.Join(parentPath, file.Name())
		fileDto := NewFileDto(
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

func safeJoin(userInput string) (string, error) {
	fullPath := filepath.Join(config.AbsRootDir, userInput)
	cleanedPath := filepath.Clean(fullPath)
	absFullPath, err := filepath.Abs(cleanedPath)

	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	if !strings.HasPrefix(absFullPath, config.AbsRootDir) {
		return "", fmt.Errorf("path traversal attempt detected: %s", userInput)
	}

	return absFullPath, nil
}

func getDbFiles(db *gorm.DB, parentPath string) (map[string]models.File, error) {
	var files []models.File
	err := db.Where("parent_path = ?", parentPath).Find(&files).Error
	if err != nil {
		return nil, err
	}

	filesMap := make(map[string]models.File)
	for _, file := range files {
		filesMap[file.Name] = file
	}

	return filesMap, nil
}

func OpenFile(path string) (*os.File, error) {
	safePath, err := safeJoin(path)
	if err != nil {
		return nil, err
	}
	return os.Open(safePath)
}
