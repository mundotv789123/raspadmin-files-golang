package icongenerator

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/djherbis/times"
	"github.com/google/uuid"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/database"
	"github.com/mundotv789123/raspadmin/internal/models"
	"github.com/mundotv789123/raspadmin/jobs/icon_generator/generator"
	"github.com/mundotv789123/raspadmin/repository"
	"gorm.io/gorm"
)

func RunGenerator() error {
	erro := processFile(config.AbsRootDir, database.DB)
	if erro != nil {
		return erro
	}
	return nil
}

func processFile(path string, db *gorm.DB) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	parentPath := path[len(config.AbsRootDir):]
	filesDb, err := repository.GetFilesMapFromParentPath(db, parentPath)
	dirs := make([]string, 0)
	for _, file := range files {
		filePath := filepath.Join(parentPath, file.Name())
		fullPath := filepath.Join(config.AbsRootDir, filePath)

		if file.IsDir() {
			dirs = append(dirs, filePath)
			err := processFile(fullPath, db)

			if err != nil {
				return err
			}
			continue
		}

		fileEntity, exists := filesDb[file.Name()]
		if !exists {
			fileEntity = *models.NewFile(file.Name(), filePath, &parentPath)
		} else {
			delete(filesDb, file.Name())
		}

		err := db.Save(&fileEntity).Error
		if err != nil {
			return err
		}

		contentType := mime.TypeByExtension(filepath.Ext(file.Name()))
		gen, ok := generator.GetGenerator(contentType)
		if !ok {
			continue
		}

		ok, err = doGenerateIcon(&fileEntity, fullPath, db)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}
		if fileEntity.IconPath == nil {
			iconPath := fmt.Sprintf("%s/_%s.jpg", config.CacheDir, uuid.New().String())
			fileEntity.IconPath = &iconPath
		}

		iconFullPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
		ok, err = generator.GenerateIcon(fullPath, iconFullPath, gen)

		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		fileEntity.SetIconPath(fileEntity.IconPath)

		err = db.Save(&fileEntity).Error
		if err != nil {
			return err
		}
	}
	for _, fileEntity := range filesDb {
		if fileEntity.IconPath != nil {
			fileIconPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
			err := os.Remove(fileIconPath)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}
			}
		}
		db.Delete(fileEntity)
	}
	return nil
}

func doGenerateIcon(fileEntity *models.File, fullPath string, db *gorm.DB) (bool, error) {
	if fileEntity.GenerateIcon {
		return true, nil
	}
	var createdAt int64
	var updatedAt int64

	t, err := times.Stat(fullPath)
	if err != nil {
		return false, err
	}
	createdAt = int64(t.BirthTime().Unix())
	updatedAt = int64(t.ModTime().Unix())

	info, err := os.Stat(filepath.Join(config.AbsRootDir, *fileEntity.IconPath))
	if errors.Is(err, os.ErrNotExist) || info == nil {
		fileEntity.IconPath = nil
	}

	if fileEntity.IconPath != nil && fileEntity.CreatedAtUnix == createdAt && fileEntity.UpdatedAtUnix == updatedAt {
		return false, nil
	}

	fileEntity.CreatedAtUnix = createdAt
	fileEntity.UpdatedAtUnix = updatedAt
	fileEntity.SetGenerateIcon()
	err = db.Save(&fileEntity).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
