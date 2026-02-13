package database

import (
	"github.com/glebarez/sqlite"
	"github.com/mundotv789123/raspadmin/internal/database/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDbConnection() (*gorm.DB, error) {
	sqliteConn := sqlite.Open("database.db")
	db, err := gorm.Open(sqliteConn, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	migrations := []interface{}{
		&models.File{},
		&models.Diretory{},
		&models.User{},
		&models.UserSession{},
	}

	for _, model := range migrations {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	return nil
}
