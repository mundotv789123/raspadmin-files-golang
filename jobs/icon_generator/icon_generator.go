package icongenerator

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"regexp"

	"github.com/djherbis/times"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/database"
	"github.com/mundotv789123/raspadmin/internal/models"
	"github.com/mundotv789123/raspadmin/jobs/icon_generator/generator"
	"github.com/mundotv789123/raspadmin/repository"
	"gorm.io/gorm"
)

var (
	videoTypeRegex = regexp.MustCompile(`^video/(mp4|mkv|webm)$`)
	audioTypeRegex = regexp.MustCompile(`^audio/(mpeg)$`)
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
			println("Processado diretório:", filePath)
			if err != nil {
				return err
			}
			continue
		}

		fileEntity, exists := filesDb[file.Name()]
		if !exists {
			fileEntity = *models.NewFile(file.Name(), filePath, &parentPath)
		}
		err := db.Save(&fileEntity).Error
		if err != nil {
			return err
		}

		contentType := mime.TypeByExtension(filepath.Ext(file.Name()))
		generator, ok := getGenerator(contentType)
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
		iconPath := "" //TODO: definir extensão correta, se já existir uma imagem no banco, sobescreva
		err = generateIcon(fullPath, iconPath, generator)
		if err != nil {
			return err
		}
		fileEntity.SetIconPath(&iconPath)

		// TODO: salvar dados após confirmar geração de icone
		// err = db.Save(&fileEntity).Error
		// if err != nil {
		// 	return err
		// }
		delete(filesDb, file.Name())
	}
	for _, fileEntity := range filesDb {
		fmt.Printf("Arquivo removido do sistema: %s\n", fileEntity.FilePath)
		// TODO: implementar remoção do arquivo do banco de dados
		// ignorar pasta de cache
		// apagar icone da pasta de cache
	}

	// varrer diretorios dirs e verificar se no banca parentDirs que n foram encontrados, e apaga-los do banco
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

	if fileEntity.CreatedAtUnix == createdAt && fileEntity.UpdatedAtUnix == updatedAt {
		return false, nil
	}
	fileEntity.CreatedAtUnix = createdAt
	fileEntity.UpdatedAtUnix = updatedAt
	err = db.Save(&fileEntity).Error
	if err != nil {
		return false, err
	}
	return false, nil
}

func getGenerator(contentType string) (IconGenerator, bool) {
	if videoTypeRegex.MatchString(contentType) {
		return &generator.IconVideoGenerator{}, true
	}
	if audioTypeRegex.MatchString(contentType) {
		return &generator.IconAudioGenerator{}, true
	}
	return nil, false

}

type IconGenerator interface {
	Generate(filePath string, iconPath string) error
}

func generateIcon(filePath string, iconPath string, generator IconGenerator) error {
	return generator.Generate(filePath, iconPath)
}
