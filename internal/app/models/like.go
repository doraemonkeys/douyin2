package models

import (
	"sync"
	"time"
)

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

// UserLikeCache 存放在sync.Map中, key为videoID
// type UserLikeCache struct {
// 	// 关注列表本来应该在前端缓存，但是青训营的API要求了，所以这里也缓存一份
// 	IsFollowed bool
// 	VideoCache VideoCacheModel
// }

type UserLike_VideoAndAuthor struct {
	VideoID  uint
	AuthorID uint
	// 视频的作者是否被此用户关注，而不是查询此用户的人是否关注了作者
	IsFollowed bool
}

type UserLikeCacheModel struct {
	// map[videoID]UserLike_VideoAndAuthor
	VideoIDMap sync.Map
}
