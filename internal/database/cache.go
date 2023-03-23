package database

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/pkg/cache"
)

var videoInfoCacher cache.Cacher[uint, models.VideoCacheModel]
var videoCacheInitOnce sync.Once

func InitVideoInfoCacher(cap int) {
	videoCacheInitOnce.Do(func() {
		videoInfoCacher = cache.NewPriorCircularMap[uint, models.VideoCacheModel](cap)
	})
}

// GetVideoInfoCacher
// 获取视频信息缓存器,请不要写入结构体中的指针或引用类型，这是并发不安全的
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

func GetUserCacher() cache.Cacher[uint, models.UserCacheModel] {
	return userCacher
}
