package models

import "time"

type UserSession struct {
	Id uint `gorm:"primaryKey;column:id"`

	RefreshToken string `gorm:"column:refresh_token"`
	ExpireAt     uint64 `gorm:"column:expire_at"`

	createdAt uint64 `gorm:"column:created_at;autoCreateTime"`
	updatedAt uint64 `gorm:"column:updated_at;autoUpdateTime"`
}

func (u *UserSession) SetRefreshToken(refreshToken string, expireMinutes int) {
	u.RefreshToken = refreshToken

	u.ExpireAt = uint64(time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix())
	u.updatedAt = uint64(time.Now().Unix())
}
