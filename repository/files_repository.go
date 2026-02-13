package repository

import (
	"github.com/mundotv789123/raspadmin/internal/models"
	"gorm.io/gorm"
)

func GetFilesMapFromParentPath(db *gorm.DB, parentPath string) (map[string]models.File, error) {
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
