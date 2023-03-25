package response

import (
	"math"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/sirupsen/logrus"
)

type VideoListResponse struct {
	CommonResponse
	// 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time参数
	NextTime  int         `json:"next_time"`
	VideoList []VideoList `json:"video_list"`
}

// SetValues 设置值
// videoList 中的视频需包含作者信息
func (v *VideoListResponse) SetValues(videoList []models.VideoModel, QueryerLiks map[uint]bool, FollowedMap map[uint]bool) {
	var nextTime int = math.MaxInt
	for _, val := range videoList {
		var video VideoList
		video.ID = int(val.ID)
		video.PlayURL = val.URL
		video.CoverURL = val.CoverURL
		video.FavoriteCount = int(val.LikeCount)
		video.CommentCount = int(val.CommentCount)
		video.Title = val.Title
		video.IsFavorite = QueryerLiks[val.ID]
		video.Author.SetValue(val.Author, FollowedMap[val.AuthorID])
		if int(val.CreatedAt.UnixMilli()) < nextTime {
			nextTime = int(val.CreatedAt.UnixMilli())
		}
		v.VideoList = append(v.VideoList, video)
	}
	v.NextTime = nextTime
	logrus.Trace("nextTime:", nextTime)
	logrus.Trace("videoList:", v.VideoList)
}
