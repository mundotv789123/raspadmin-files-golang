package icongenerator

import (
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/djherbis/times"
	"github.com/google/uuid"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/database"
	"github.com/mundotv789123/raspadmin/internal/models"
	"github.com/mundotv789123/raspadmin/internal/validations"
	"github.com/mundotv789123/raspadmin/jobs/icon_generator/generator"
	"github.com/mundotv789123/raspadmin/jobs/icon_generator/watch"
	"github.com/mundotv789123/raspadmin/repository"
	"gorm.io/gorm"
)

var (
	fileLocks   = make(map[string]*sync.Mutex)
	fileLocksMu sync.Mutex
)

func LockFile(path string) func() {
	fileLocksMu.Lock()
	mtx, ok := fileLocks[path]
	if !ok {
		mtx = &sync.Mutex{}
		fileLocks[path] = mtx
	}
	fileLocksMu.Unlock()

	mtx.Lock()
	return func() {
		mtx.Unlock()
	}
}

func StartBackgroundService() {
	slog.Info("starting icon generator background service")

	go func() {
		if err := RunGenerator(); err != nil {
			slog.Error(fmt.Sprintf("error running initial icon generator: %v", err))
		}
	}()

	go func() {
		watch.InitWatch(config.AbsRootDir, func(filePath string) {
			if err := ProcessSingleFile(filePath); err != nil {
				slog.Error(fmt.Sprintf("error processing watched file %s: %v", filePath, err))
			}
		})
	}()
}

func RunGenerator() error {
	slog.Info("starting files process")
	if err := processFile(config.AbsRootDir, database.DB); err != nil {
		return err
	}
	slog.Info("files process finished")
	return nil
}

func ProcessSingleFile(fullPath string) error {
	fullPath = filepath.Clean(fullPath)

	if strings.HasPrefix(fullPath, config.CacheDirAds) {
		return nil
	}
	if validations.IsTempOrHiddenFile(filepath.Base(fullPath)) {
		return nil
	}
	if !strings.HasPrefix(fullPath, config.AbsRootDir) {
		return nil
	}

	unlock := LockFile(fullPath)
	defer unlock()

	db := database.DB
	if db == nil {
		return errors.New("database connection not initialized")
	}

	relPath := fullPath[len(config.AbsRootDir):]
	if relPath == "" {
		relPath = "/"
	}

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		// File was removed
		fileEntity, err := repository.GetFileFromPath(db, relPath)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		if fileEntity.IconPath != nil && *fileEntity.IconPath != "" {
			iconFullPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
			if err := os.Remove(iconFullPath); err != nil && !errors.Is(err, os.ErrNotExist) {
				slog.Error(fmt.Sprintf("error remove icon file %s: %v", iconFullPath, err))
			}
		}
		slog.Info(fmt.Sprintf("file deleted from database %s", fileEntity.FilePath))
		return db.Delete(fileEntity).Error
	}
	if err != nil {
		return err
	}

	if info.IsDir() {
		return processFile(fullPath, db)
	}

	parentPath := filepath.Dir(relPath)
	if parentPath == "." {
		parentPath = ""
	}
	fileName := info.Name()

	fileEntity, err := repository.GetFileFromPath(db, relPath)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fileEntity = models.NewFile(fileName, relPath, &parentPath)
			slog.Info(fmt.Sprintf("file %s will be created in database.", relPath))
		} else {
			return err
		}
	}

	if err := db.Save(fileEntity).Error; err != nil {
		return fmt.Errorf("error save file in db: %w", err)
	}

	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	gen, ok := generator.GetGenerator(contentType)
	if !ok {
		slog.Debug(fmt.Sprintf("no generator found to file %s.", contentType))
		return nil
	}

	ok, err = doGenerateIcon(fileEntity, fullPath, db)
	if err != nil {
		return fmt.Errorf("error save file doGenerateIcon: %w", err)
	}
	if !ok {
		return nil
	}

	if fileEntity.IconPath == nil {
		iconPath := fmt.Sprintf("%s/_%s.jpg", config.CacheDir, uuid.New().String())
		fileEntity.IconPath = &iconPath
	}

	iconFullPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
	slog.Info(fmt.Sprintf("generating icon to file %s saving in %s", fullPath, iconFullPath))
	ok, err = generator.GenerateIcon(fullPath, iconFullPath, gen)
	if err != nil {
		return fmt.Errorf("error save file generate icon: %w", err)
	}

	if ok {
		fileEntity.SetIconPath(fileEntity.IconPath)
	} else {
		slog.Info(fmt.Sprintf("icon %s was not generated", iconFullPath))
		fileEntity.SetIconPath(nil)
	}

	return db.Save(fileEntity).Error
}

func processFile(path string, db *gorm.DB) error {
	slog.Debug(fmt.Sprintf("reading dir %s", path))
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error process file: %s, %s", err, path)
	}

	parentPath := path[len(config.AbsRootDir):]
	filesDb, err := repository.GetFilesMapFromParentPath(db, parentPath)
	slog.Debug(fmt.Sprintf("%d file(s) were found in the database.", len(filesDb)))

	for _, file := range files {
		filePath := filepath.Join(parentPath, file.Name())
		fullPath := filepath.Join(config.AbsRootDir, filePath)

		if file.IsDir() {
			if strings.HasPrefix(fullPath, config.CacheDirAds) || validations.IsTempOrHiddenFile(file.Name()) {
				continue
			}
			err := processFile(fullPath, db)

			if err != nil {
				return err
			}
			continue
		}

		unlock := LockFile(fullPath)

		fileEntity, exists := filesDb[file.Name()]
		if !exists {
			fileEntity = *models.NewFile(file.Name(), filePath, &parentPath)
			slog.Info(fmt.Sprintf("file %s will be created in the database.", filePath))
		} else {
			delete(filesDb, file.Name())
			slog.Debug(fmt.Sprintf("file %s already exists.", filePath))
		}

		if err := db.Save(&fileEntity).Error; err != nil {
			unlock()
			return fmt.Errorf("error save file in db: %s (%s, %s)", err, file.Name(), path)
		}

		contentType := mime.TypeByExtension(filepath.Ext(file.Name()))
		gen, ok := generator.GetGenerator(contentType)
		if !ok {
			slog.Debug(fmt.Sprintf("no generator found to file %s.", contentType))
			unlock()
			continue
		}

		ok, err = doGenerateIcon(&fileEntity, fullPath, db)
		if err != nil {
			unlock()
			return fmt.Errorf("error save file doGenerateIcon: %s (%s, %s)", err, file.Name(), path)
		}

		if !ok {
			unlock()
			continue
		}
		if fileEntity.IconPath == nil {
			iconPath := fmt.Sprintf("%s/_%s.jpg", config.CacheDir, uuid.New().String())
			fileEntity.IconPath = &iconPath
		}

		iconFullPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
		slog.Info(fmt.Sprintf("generating icon to file %s saving in %s", fullPath, iconFullPath))
		ok, err = generator.GenerateIcon(fullPath, iconFullPath, gen)

		if err != nil {
			unlock()
			return fmt.Errorf("error save file generate icon: %s, (%s, %s)", err, fullPath, path)
		}

		if ok {
			fileEntity.SetIconPath(fileEntity.IconPath)
		} else {
			slog.Info(fmt.Sprintf("icon %s was not generated", iconFullPath))
			fileEntity.SetIconPath(nil)
		}

		if err = db.Save(&fileEntity).Error; err != nil {
			unlock()
			return fmt.Errorf("error save file in db %s, %s", file.Name(), path)
		}

		unlock()
	}
	for _, fileEntity := range filesDb {
		if fileEntity.IconPath != nil && *fileEntity.IconPath != "" {
			fileIconPath := filepath.Join(config.AbsRootDir, *fileEntity.IconPath)
			slog.Info(fmt.Sprintf("delete icon from cache %s", fileIconPath))
			err := os.Remove(fileIconPath)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("error remove file %s", fileIconPath)
				}
			}
		}
		slog.Info(fmt.Sprintf("file deleted from database %s", fileEntity.FilePath))
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

