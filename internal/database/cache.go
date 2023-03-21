package database

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/pkg/cache"
)

var videoInfoCacher cache.Cacher[uint, models.VideoCacheModel]
var videoCacheInitOnce sync.Once

// GetVideoInfoCacher
// 获取视频信息缓存器,请不要写入结构体中的指针或引用类型，这是并发不安全的
func GetVideoInfoCacher() cache.Cacher[uint, models.VideoCacheModel] {
	videoCacheInitOnce.Do(func() {

	})
	return videoInfoCacher
}

var videoCommentCacher cache.Cacher[uint, []models.CommentModel]
var videoCommentCacheInitOnce sync.Once

func GetVideoCommentCacher() cache.Cacher[uint, []models.CommentModel] {
	videoCommentCacheInitOnce.Do(func() {

	})
	return videoCommentCacher
}

var userCacher cache.Cacher[uint, models.UserCacheModel]
var userCacheInitOnce sync.Once

// GetUserCacher
// 获取用户信息缓存器,请不要写入结构体中的指针或引用类型，这是并发不安全的
func GetUserCacher() cache.Cacher[uint, models.UserCacheModel] {
	userCacheInitOnce.Do(func() {

	})
	return userCacher
}
