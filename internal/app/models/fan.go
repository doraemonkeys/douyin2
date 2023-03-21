package models

import "time"

// 粉丝列表
type UserFansModel struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	FanID     uint
	CreatedAt time.Time
}
