package models

import "time"

const (
	UserFollowerModelTableName        = "user_follower_models"
	UserFollowerModelTable_UserID     = "user_id"
	UserFollowerModelTable_FollowerID = "follower_id"
)

// 关注列表
type UserFollowerModel struct {
	CreatedAt  time.Time
	UserID     uint `gorm:"primary_key"`
	FollowerID uint `gorm:"primary_key"`
}
