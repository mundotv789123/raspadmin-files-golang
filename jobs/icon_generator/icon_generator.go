package icongenerator

import (
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"strings"

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
	slog.Info("starting files process")
	if err := processFile(config.AbsRootDir, database.DB); err != nil {
		return err
	}
	slog.Info("files process finished")
	return nil
}

func processFile(path string, db *gorm.DB) error {
	slog.Debug("reading dir %s", path)
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error process file: %s, %s", err, path)
	}

	parentPath := path[len(config.AbsRootDir):]
	filesDb, err := repository.GetFilesMapFromParentPath(db, parentPath)
	slog.Debug("%d file(s) were found in the database.", len(filesDb))

	for _, file := range files {
		filePath := filepath.Join(parentPath, file.Name())
		fullPath := filepath.Join(config.AbsRootDir, filePath)

		if file.IsDir() {
			if strings.HasPrefix(fullPath, config.CacheDirAds) {
				continue
			}
			err := processFile(fullPath, db)

			if err != nil {
				return err
			}
			continue
		}

		fileEntity, exists := filesDb[file.Name()]
		if !exists {
			fileEntity = *models.NewFile(file.Name(), filePath, &parentPath)
			slog.Info("file %s will be created in the database.", filePath)
		} else {
			delete(filesDb, file.Name())
			slog.Debug("file %s already exists.", filePath) //
		}

		if err := db.Save(&fileEntity).Error; err != nil {
			return fmt.Errorf("error save file in db: %s (%s, %s)", err, file.Name(), path)
		}

		contentType := mime.TypeByExtension(filepath.Ext(file.Name()))
		gen, ok := generator.GetGenerator(contentType)
		if !ok {
			slog.Debug("no generator found to file %s.", contentType)
			continue
		}

		ok, err = doGenerateIcon(&fileEntity, fullPath, db)
		if err != nil {
			return fmt.Errorf("error save file doGenerateIcon: %s (%s, %s)", err, file.Name(), path)
		}

		if !ok {
			continue
		}
		if fileEntity.IconPath == nil {
			iconPath := fmt.Sprintf("%s/_%s.jpg", config.CacheDir, uuid.New().String())
			fileEntity.IconPath = &iconPath
		}

		iconFullPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
		slog.Info("generating icon to file %s saving in %s", fullPath, iconFullPath)
		ok, err = generator.GenerateIcon(fullPath, iconFullPath, gen)

		if err != nil {
			return fmt.Errorf("error save file generate icon: %s, (%s, %s)", err, fullPath, path)
		}

		if ok {
			fileEntity.SetIconPath(fileEntity.IconPath)
		} else {
			slog.Info("icon %s was not generated", iconFullPath)
			fileEntity.SetIconPath(nil)
		}

		if err = db.Save(&fileEntity).Error; err != nil {
			return fmt.Errorf("error save file in db %s, %s", file.Name(), path)
		}
	}
	for _, fileEntity := range filesDb {
		if fileEntity.IconPath != nil && *fileEntity.IconPath != "" {
			fileIconPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
			slog.Info("delete icon from cache %s", fileIconPath)
			err := os.Remove(fileIconPath)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("error remove file %s", fileIconPath)
				}
			}
		}
		slog.Info("file deleted from database %s", fileEntity.FilePath)
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

	if fileEntity.IconPath != nil {
		info, err := os.Stat(filepath.Join(config.AbsRootDir, *fileEntity.IconPath))
		if errors.Is(err, os.ErrNotExist) || info == nil {
			fileEntity.IconPath = nil
		}
	}

	if fileEntity.IconPath != nil && fileEntity.CreatedAtUnix == createdAt && fileEntity.UpdatedAtUnix == updatedAt {
		return false, nil
	}

	fileEntity.CreatedAtUnix = createdAt
	fileEntity.UpdatedAtUnix = updatedAt
	if err := fileEntity.SetGenerateIcon(); err != nil {
		return false, err
	}

	if err = db.Save(&fileEntity).Error; err != nil {
		return false, err
	}
	return true, nil
}
