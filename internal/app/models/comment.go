package models

import (
	"sync"

	"gorm.io/gorm"
)

const CommentModelTableName = "comment_models"

const (
	CommentModelTable_VideoID     = "video_id"
	CommentModelTable_UserID      = "user_id"
	CommentModelTable_Content     = "content"
	CommentModelPreload_Commenter = "Commenter"
)

// 评论列表
type CommentModel struct {
	gorm.Model
	VideoID uint
	UserID  uint
	Content string `gorm:"type:text"`
	Video   VideoModel
	// 使用前需确保里面有数据
	Commenter UserModel `gorm:"foreignKey:UserID"`
}

type CommentCacheModel struct {
	//map[commentID]CommentModel
	CacheMap map[uint]CommentModel
	MapLock  *sync.RWMutex
}

func NewCommentCacheModel() CommentCacheModel {
	return CommentCacheModel{
		CacheMap: make(map[uint]CommentModel),
		MapLock:  new(sync.RWMutex),
	}
}
