package models

import "time"

const (
	UserLikeModelTableName       = "user_like_models"
	UserLikeModelTable_UserID    = "user_id"
	UserLikeModelTable_VideoID   = "video_id"
	UserLikeModelTable_CreatedAt = "created_at"
)

// 点赞列表
type UserLikeModel struct {
	UserID    uint `gorm:"primarykey"`
	VideoID   uint `gorm:"primarykey"`
	CreatedAt time.Time
}
