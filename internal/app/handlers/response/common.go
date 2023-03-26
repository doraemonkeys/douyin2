package response

import (
	"net/http"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/gin-gonic/gin"
)

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

// status code
const (
	//Success      = 200
	Success      = 0
	Failed       = 500
	TokenExpired = 401
)

// errors
const (
	ErrServerInternal = "服务器内部错误"
	ErrUserNotLogin   = "用户未登录"
	ErrUserNotExists  = "用户不存在"
	ErrUserExists     = "用户已存在"
	ErrUserPassword   = "密码错误"
	ErrUserToken      = "用户token错误"
	ErrUserTokenExp   = "用户token过期"
	ErrInvalidParams  = "参数错误"
	ErrDBEmpty        = "已经没有更多视频了"
)

// success response
const (
	SuccessMsg      = "success"
	QuerySuccessMsg = "查询成功"
	EmptyVideoList  = "暂无视频"
)

type CommonResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// 通用响应错误
func ResponseError(c *gin.Context, errMsg string) {
	var res CommonResponse
	res.StatusCode = Failed
	res.StatusMsg = errMsg
	c.JSON(http.StatusOK, res)
}

// 通用响应成功
func ResponseSuccess(c *gin.Context, msg string) {
	var res CommonResponse
	res.StatusCode = Success
	res.StatusMsg = msg
	c.JSON(http.StatusOK, res)
}

// isFollow 查询者是否关注了作者
func (a *Author) SetValue(user models.UserModel, isFollow bool) {
	a.ID = int(user.ID)
	a.Name = user.Username
	a.FollowCount = int(user.FollowerCount)
	a.FollowerCount = int(user.FanCount)
	a.IsFollow = isFollow
}

// func (a *Author) SetValueFromCache(cache models.UserLikeCache) {
// 	var user models.UserModel
// 	user.SetValueFromCacheModel(cache.VideoCache.Author)
// 	a.SetValue(user, cache.IsFollowed)
// }

// func (v *VideoList) SetValueFromUserLikeCache(cache models.UserLikeCache) {
// 	v.ID = int(cache.VideoCache.AuthorID)
// 	v.PlayURL = cache.VideoCache.URL
// 	v.CoverURL = cache.VideoCache.CoverURL
// 	v.FavoriteCount = int(cache.VideoCache.LikeCount)
// 	v.CommentCount = int(cache.VideoCache.CommentCount)
// 	v.Title = cache.VideoCache.Title
// 	v.Author.SetValueFromCache(cache)
// }
