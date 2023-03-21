package models

import "gorm.io/gorm"

const CommentModelTableName = "comment_models"

// 评论列表
type CommentModel struct {
	gorm.Model
	VideoID uint
	UserID  uint
	Content string `gorm:"type:text"`
	Video   VideoModel
	User    UserModel
}
