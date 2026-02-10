package models

import "time"

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

	UpdatedAt     *time.Time `gorm:"-"`
	UpdatedAtUnix int64      `gorm:"column:updated_at;autoCreateTime"`
	CreatedAt     *time.Time `gorm:"-"`
	CreatedAtUnix int64      `gorm:"column:created_at;autoCreateTime"`
}

func NewFile(name string, filePath string, parentPath *string) *File {
	return &File{
		Name:         name,
		FilePath:     filePath,
		ParentPath:   parentPath,
		GenerateIcon: false,
	}
}

func SetGenerateIcon(file *File) {
	file.GenerateIcon = true
	file.IconPath = nil
}

func SetIconPath(file *File, iconPath *string) {
	file.IconPath = iconPath
	file.GenerateIcon = false
}
