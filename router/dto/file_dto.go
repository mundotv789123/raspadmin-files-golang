package dto

import (
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/times"
)

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
		Path:      path,
		Type:      contentType,
		Icon:      iconPath,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Open:      open,
	}
}
