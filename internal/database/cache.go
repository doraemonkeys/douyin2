package database

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/pkg/cache"
)

// 青训营API不合理的地方：
// 1. 前端不缓存关注列表，每次都要从后端获取，导致每次获取关注信息都要从MySQL中查询
// 2. 返回的视频列表(feed流，喜欢列表，评论列表)中要求包含作者的动态信息(如粉丝总数)。
//    实际应用中，用户不可能点进每个视频作者或评论用户的主页。

var videoInfoCacher cache.Cacher[uint, models.VideoCacheModel]
var videoCacheInitOnce sync.Once

func InitVideoInfoCacher(cap int) {
	videoCacheInitOnce.Do(func() {
		videoInfoCacher = cache.NewARC[uint, models.VideoCacheModel](cap)
	})
}

// GetVideoInfoCacher
// 获取视频信息缓存器,请不要写入结构体中的指针或引用类型，这是并发不安全的。
// 请注意Author字段可能有值，也可能没有值。
func GetVideoInfoCacher() cache.Cacher[uint, models.VideoCacheModel] {
	return videoInfoCacher
}

var videoCommentCacher cache.Cacher[uint, models.CommentModel]
var videoCommentCacherInitOnce sync.Once

func InitVideoCommentCacher(cap int) {
	videoCommentCacherInitOnce.Do(func() {
		videoCommentCacher = cache.NewARC[uint, models.CommentModel](cap)
	})
}

func GetVideoCommentCacher() cache.Cacher[uint, models.CommentModel] {
	return videoCommentCacher
}

var userCacher cache.Cacher[uint, models.UserCacheModel]
var userCacheInitOnce sync.Once

func InitUserCacher(cap int) {
	userCacheInitOnce.Do(func() {
		userCacher = cache.NewARC[uint, models.UserCacheModel](cap)
	})
}

func GetUserInfoCacher() cache.Cacher[uint, models.UserCacheModel] {
	return userCacher
}

var userFavoriteCacher cache.Cacher[uint, models.UserLikeCacheModel]
var userFavoriteCacherInitOnce sync.Once

func InitUserFavoriteCacher(cap int) {
	userFavoriteCacherInitOnce.Do(func() {
		userFavoriteCacher = cache.NewARC[uint, models.UserLikeCacheModel](cap)
	})
}

func GetUserFavoriteCacher() cache.Cacher[uint, models.UserLikeCacheModel] {
	return userFavoriteCacher
}
