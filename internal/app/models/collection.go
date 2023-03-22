package models

import (
	"time"
)

// 收藏列表
type UserCollectionModel struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	VideoID   uint
	CreatedAt time.Time
}
