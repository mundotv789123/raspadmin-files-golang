package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mundotv789123/raspadmin/internal/config"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Path  string `json:"path"`
	Open  bool   `json:"open,omitempty"`
}

type FileContent interface {
	Name() string
	IsDir() bool
}

func NewFileInfo(file FileContent, path string, open bool) FileInfo {
	return FileInfo{
		Name:  file.Name(),
		IsDir: file.IsDir(),
		Path:  path,
	}
}

var (
	ErrFileNotFound = errors.New("File or directory not found")
)

func GetFiles(path string) ([]FileInfo, error) {
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
		return []FileInfo{NewFileInfo(fileState, path, true)}, nil
	}

	files, err := os.ReadDir(path)
	filesList := make([]FileInfo, len(files))
	for i, file := range files {
		filePath := filepath.Join(path, file.Name())[len(config.AbsRootDir):]
		filesList[i] = NewFileInfo(file, filePath, false)
	}
	return filesList, nil
}

func safeJoin(userInput string) (string, error) {
	fullPath := filepath.Join(config.RootDir, userInput)
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

func OpenFile(path string) (*os.File, error) {
	safePath, err := safeJoin(path)
	if err != nil {
		return nil, err
	}
	return os.Open(safePath)
}
