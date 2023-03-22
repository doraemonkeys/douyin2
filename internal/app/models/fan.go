package models

import "time"

// 粉丝列表
type UserFanModel struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	FanID     uint
	CreatedAt time.Time
}
