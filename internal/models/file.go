package models

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mundotv789123/raspadmin/internal/config"
)

type File struct {
	Id uint `gorm:"primaryKey;column:id"`

	Name string  `gorm:"not null;column:name"`
	Size *int64  `gorm:"column:size"`
	Dir  bool    `gorm:"not null;column:is_dir"`
	Type *string `gorm:"column:type"`

	GenerateIcon bool    `gorm:"not null;column:generate_icon;default:false"`
	IconPath     *string `gorm:"column:icon_path"`

	FilePath   string  `gorm:"not null;column:file_path;uniqueIndex"`
	ParentPath *string `gorm:"column:parent_path"`

	UpdatedAtUnix int64 `gorm:"column:updated_at;autoCreateTime"`
	CreatedAtUnix int64 `gorm:"column:created_at;autoCreateTime"`
}

func NewFile(name string, filePath string, parentPath *string) *File {
	return &File{
		Name:         name,
		FilePath:     filePath,
		ParentPath:   parentPath,
		GenerateIcon: false,
	}
}

func (file *File) SetGenerateIcon() error {
	if file.IconPath != nil && *file.IconPath != "" {
		fileIconPath := filepath.Join(config.AbsRootDir, *file.IconPath)
		slog.Info(fmt.Sprintf("delete icon from cache %s", fileIconPath))
		err := os.Remove(fileIconPath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	file.GenerateIcon = true
	file.IconPath = nil
	return nil
}

func (file *File) SetIconPath(iconPath *string) {
	file.IconPath = iconPath
	file.GenerateIcon = false
}
