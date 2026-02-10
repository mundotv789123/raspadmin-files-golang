package models

type User struct {
	Id       uint   `gorm:"primaryKey;column:id"`
	Username string `gorm:"not null;uniqueIndex;column:username"`
	Password string `gorm:"not null;column:password"`
}
