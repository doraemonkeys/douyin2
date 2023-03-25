package response

import "github.com/Doraemonkeys/douyin2/internal/app/models"

type QueryFavorVideoListResponse struct {
	CommonResponse
	VideoList []VideoList `json:"video_list"`
}

func (r *QueryFavorVideoListResponse) SetValues(videoList []models.VideoModel, FollowedMap map[uint]bool) {
	r.VideoList = make([]VideoList, len(videoList))
	for i, val := range videoList {
		var video VideoList
		video.ID = int(val.ID)
		video.PlayURL = val.URL
		video.CoverURL = val.CoverURL
		video.FavoriteCount = int(val.LikeCount)
		video.CommentCount = int(val.CommentCount)
		video.Title = val.Title
		video.Author.SetValue(val.Author, FollowedMap[val.AuthorID])
		r.VideoList[i] = video
	}
}

// func (r *QueryFavorVideoListResponse) SetValuesFromCache(likes []models.UserLikeCache) {
// 	r.VideoList = make([]VideoList, len(likes))
// 	for _, val := range likes {
// 		var video VideoList
// 		video.SetValueFromUserLikeCache(val)
// 	}
// }
