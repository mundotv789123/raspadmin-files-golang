package database

import (
	"github.com/mundotv789123/raspadmin/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDbConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
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
	err := db.AutoMigrate(&models.File{})
	if err != nil {
		return err
	}
	return nil
}
