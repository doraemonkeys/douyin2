package response

import "github.com/Doraemonkeys/douyin2/internal/app/models"

type PublishVedioResponse struct {
	CommonResponse
}

type QueryVideoListResponse struct {
	CommonResponse
	VideoList []VideoList2 `json:"video_list"`
}

// QueryVideoListResponse
type Author2 struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// QueryVideoListResponse
type VideoList2 struct {
	//视频ID
	ID            int     `json:"id"`
	Author        Author2 `json:"author"`
	PlayURL       string  `json:"play_url"`
	CoverURL      string  `json:"cover_url"`
	FavoriteCount int     `json:"favorite_count"`
	CommentCount  int     `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
	Title         string  `json:"title"`
}

// 这里的VideoModel里不需要存储作者的信息，单独给出
// 这里查出了查询者的所有喜欢的视频，然后在遍历视频列表的时候，判断是否是喜欢的视频
// 以后优化
func (v *QueryVideoListResponse) SetValues(videoList []models.VideoModel, targetUser models.UserModel, QueryerLiks map[uint]struct{}, followed bool) {
	for _, val := range videoList {
		var video VideoList2
		video.ID = int(val.ID)
		video.Author.ID = int(targetUser.ID)
		video.Author.Name = targetUser.Username
		video.Author.FollowCount = int(targetUser.FollowerCount)
		video.Author.FollowerCount = int(targetUser.FanCount)
		if followed {
			video.Author.IsFollow = true
		}
		video.PlayURL = val.URL
		video.CoverURL = val.CoverURL
		video.FavoriteCount = int(val.LikeCount)
		video.CommentCount = int(val.CommentCount)
		video.Title = val.Title
		if _, ok := QueryerLiks[val.ID]; ok {
			video.IsFavorite = true
		}
		v.VideoList = append(v.VideoList, video)
	}
}
