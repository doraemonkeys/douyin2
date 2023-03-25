package models

import "gorm.io/gorm"

const VideoTitleMaxByteLength = 199
const VideoTitleMaxRuneLength = 60

const (
	VideoModelTableName              = "videos_models"
	VideoModelTable_LikeCount        = "like_count"
	VideoModelTable_CreatedAt        = "created_at"
	VideoModelTable_AuthorID         = "author_id"
	VideoModelTable_LikesSlice       = "Likes"
	VideoModelTable_CollectionsSlice = "Collections"
)

// 视频
type VideoModel struct {
	gorm.Model
	//Title 视频标题
	Title string `gorm:"size:200"`
	// 视频保存地址
	//Path string `gorm:"size:200;unique"`
	// 视频封面保存地址
	//CoverPath    string `gorm:"size:200"`
	StorageID    uint   `gorm:"unique_index"`
	URL          string `gorm:"size:200"`
	CoverURL     string `gorm:"size:200"`
	Author       UserModel
	AuthorID     uint
	LikeCount    uint
	CommentCount uint
	Comments     []CommentModel `gorm:"foreignKey:VideoID"`
	Likes        []UserModel    `gorm:"many2many:user_like;joinForeignKey:VideoID;joinReferences:UserID"`
	Collections  []UserModel    `gorm:"many2many:user_collection;joinForeignKey:VideoID;joinReferences:UserID"`
}

func (v *VideoModel) TableName() string {
	return VideoModelTableName
}

type VideoCacheModel struct {
	gorm.Model
	Title        string
	StorageID    uint
	URL          string
	CoverURL     string
	AuthorID     uint
	LikeCount    uint
	CommentCount uint
	//Author       UserCacheModel
}

func (v *VideoCacheModel) SetValue(other VideoModel) {
	v.ID = other.ID
	v.CreatedAt = other.CreatedAt
	v.UpdatedAt = other.UpdatedAt
	v.DeletedAt = other.DeletedAt
	v.Title = other.Title
	v.StorageID = other.StorageID
	v.URL = other.URL
	v.CoverURL = other.CoverURL
	v.AuthorID = other.AuthorID
	v.LikeCount = other.LikeCount
	v.CommentCount = other.CommentCount
}

func (v *VideoModel) SetValueFromCacheModel(other VideoCacheModel) {
	v.ID = other.ID
	v.CreatedAt = other.CreatedAt
	v.UpdatedAt = other.UpdatedAt
	v.DeletedAt = other.DeletedAt
	v.Title = other.Title
	v.StorageID = other.StorageID
	v.URL = other.URL
	v.CoverURL = other.CoverURL
	v.AuthorID = other.AuthorID
	v.LikeCount = other.LikeCount
	v.CommentCount = other.CommentCount
}
