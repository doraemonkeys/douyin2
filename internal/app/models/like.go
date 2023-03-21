package models

import "time"

// 点赞列表
type UserLikesModel struct {
	UserID    uint `gorm:"primarykey"`
	VideoID   uint `gorm:"primarykey"`
	CreatedAt time.Time
}
