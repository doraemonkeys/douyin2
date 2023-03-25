package response

import "github.com/Doraemonkeys/douyin2/internal/app/models"

type GetUserInfoResponse struct {
	CommonResponse
	User User `json:"user"`
}
type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// 关注列表应该在前端做缓存，不然每次都要请求服务器，太浪费资源了

type QueryCommentListResponse struct {
	CommonResponse
	CommentList []CommentList `json:"comment_list"`
}

type CommentList struct {
	ID      int    `json:"id"`
	User    User   `json:"user"`
	Content string `json:"content"`
	// 评论发布日期，格式 2006-01-02 15:04
	CreateDate string `json:"create_date"`
}

// 查询者是否关注传入的user
func (u *User) SetValue(user models.UserModel, isFollow bool) {
	u.ID = int(user.ID)
	u.Name = user.Username
	u.FollowCount = int(user.FollowerCount)
	u.FollowerCount = int(user.FanCount)
	u.IsFollow = isFollow
}
