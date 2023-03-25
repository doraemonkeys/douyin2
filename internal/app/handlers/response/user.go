package response

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
