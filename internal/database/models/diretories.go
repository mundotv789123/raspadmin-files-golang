package models

type Diretory struct {
	Id        uint   `gorm:"primaryKey;column:id"`
	Path      string `gorm:"not null;column:path;uniqueIndex"`
	IsLoading bool   `gorm:"not null;column:is_loading;default:false"`
}
