package models

import (
	"gorm.io/gorm"
)

const (
	UserModelTableName              = "users_models"
	UserModelTable_Username         = "username"
	UserModelTable_LikesSlice       = "Likes"
	UserModelTable_FollowersSlice   = "Followers"
	UserModelTable_FansSlice        = "Fans"
	UserModelTable_CollectionsSlice = "Collections"
)

const DataBaseTimeFormat = "2006-01-02 15:04:05.000"

// 用户
type UserModel struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:100"`
	Password string `gorm:"size:100"`
	Email    string `gorm:"size:100"`
	//Phone    string `gorm:"uniqueIndex;size:50"`
	//总关注数
	FollowerCount uint `gorm:"type:int"`
	//总粉丝数
	FanCount uint `gorm:"type:int"`
	//总评论数
	CommentCount uint
	//总收藏数
	CollectionsCount uint `gorm:"type:int"`
	//关注列表
	Followers []UserModel `gorm:"many2many:user_follower;joinForeignKey:UserID;joinReferences:FollowerID"`
	//粉丝列表
	Fans []UserModel `gorm:"many2many:user_fan;joinForeignKey:UserID;joinReferences:FanID"`
	//点赞列表
	Likes []VideoModel `gorm:"many2many:user_like;joinForeignKey:UserID;joinReferences:VideoID"`
	//评论列表
	Comments []CommentModel `gorm:"foreignKey:UserID"`
	//收藏列表
	Collections []VideoModel `gorm:"many2many:user_collection;joinForeignKey:UserID;joinReferences:VideoID"`
	//视频列表
	Videos []VideoModel `gorm:"foreignKey:AuthorID"`
}

func (u *UserModel) TableName() string {
	return UserModelTableName
}

type UserCacheModel struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:100"`
	Password string `gorm:"size:100"`
	Email    string `gorm:"size:100"`
	//Phone    string `gorm:"uniqueIndex;size:50"`
	//总关注数
	FollowerCount uint `gorm:"type:int"`
	//总粉丝数
	FanCount uint `gorm:"type:int"`
	//总评论数
	CommentCount uint
	//总收藏数
	CollectionsCount uint `gorm:"type:int"`
	//总视频数
	VideosCount uint `gorm:"type:int"`
}

func (u *UserCacheModel) SetValue(user UserModel) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.DeletedAt = user.DeletedAt
	u.Username = user.Username
	u.Password = user.Password
	u.Email = user.Email
	u.FollowerCount = user.FollowerCount
	u.FanCount = user.FanCount
	u.CommentCount = user.CommentCount
	u.CollectionsCount = user.CollectionsCount
}

func (u *UserModel) SetValueFromCacheModel(user UserCacheModel) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.DeletedAt = user.DeletedAt
	u.Username = user.Username
	u.Password = user.Password
	u.Email = user.Email
	u.FollowerCount = user.FollowerCount
	u.FanCount = user.FanCount
	u.CommentCount = user.CommentCount
	u.CollectionsCount = user.CollectionsCount
}
