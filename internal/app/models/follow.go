package models

import "time"

const (
	UserFollowerModelTableName        = "user_follower_models"
	UserFollowerModelTable_UserID     = "user_id"
	UserFollowerModelTable_FollowerID = "follower_id"
)

// 关注列表
type UserFollowerModel struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UserID     uint
	FollowerID uint
}
