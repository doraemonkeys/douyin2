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

type Author struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type VideoList struct {
	ID            int    `json:"id"`
	Author        Author `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

func (a *Author) SetValues(user models.UserModel) {
	a.ID = int(user.ID)
	a.Name = user.Username
	a.FollowCount = int(user.FollowerCount)
	a.FollowerCount = int(user.FanCount)
}

func (v *VideoListResponse) SetValues(videoList []models.VideoModel, UsersMap map[uint]models.UserModel) {
	var nextTime int = math.MaxInt
	for _, val := range videoList {
		var video VideoList
		video.ID = int(val.ID)
		video.PlayURL = val.URL
		video.CoverURL = val.CoverURL
		video.FavoriteCount = int(val.LikeCount)
		video.CommentCount = int(val.CommentCount)
		video.Title = val.Title
		video.Author.SetValues(UsersMap[val.AuthorID])
		if int(val.CreatedAt.Unix()) < nextTime {
			nextTime = int(val.CreatedAt.Unix())
		}
		v.VideoList = append(v.VideoList, video)
	}
	logrus.Trace("nextTime:", nextTime)
	logrus.Trace("videoList:", v.VideoList)
}
