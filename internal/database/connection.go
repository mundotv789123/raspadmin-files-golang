package database

import (
	"github.com/glebarez/sqlite"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDbConnection() (*gorm.DB, error) {
	sqliteConn := sqlite.Open(config.DbFile)
	db, err := gorm.Open(sqliteConn, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Optimize SQLite performance for concurrent writes
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	db.Exec("PRAGMA busy_timeout=5000;")
	db.Exec("PRAGMA cache_size=-64000;")

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}


func runMigrations(db *gorm.DB) error {
	migrations := []any{
		&models.File{},
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
