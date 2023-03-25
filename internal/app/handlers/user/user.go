package user

import (
	"strconv"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/gin-gonic/gin"
)

// GetUserInfoDTO 获取用户信息DTO
type GetUserInfoDTO struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

const (
	GetUserInfoDTO_UserID = "user_id"
)

// GetUserInfoHandler 获取用户信息
func GetUserInfoHandler(c *gin.Context) {
	var getUserInfoDTO GetUserInfoDTO
	var res response.GetUserInfoResponse

	//获取请求参数
	id, _ := strconv.ParseUint(c.Query(GetUserInfoDTO_UserID), 10, 64)
	getUserInfoDTO.UserID = uint(id)
	if getUserInfoDTO.UserID == 0 {
		res.StatusCode = response.Failed
		res.StatusMsg = response.ErrInvalidParams
		c.JSON(200, res)
		return
	}

	temp, _ := c.Get("user")
	user := temp.(app.User)
	res.User.ID = int(user.ID)
	res.User.Name = user.Username

	//获取指定用户信息
	UserModel, err := services.GetUserById(getUserInfoDTO.UserID)
	if err != nil {
		res.StatusCode = response.Failed
		if err.Error() == response.ErrUserNotExists {
			res.StatusMsg = response.ErrUserNotExists
		} else {
			res.StatusMsg = response.ErrServerInternal
		}
		c.JSON(200, res)
		return
	}
	res.User.FollowCount = int(UserModel.FollowerCount)
	res.User.FollowerCount = int(UserModel.FanCount)

	//判断是否关注
	res.User.IsFollow = services.QueryUserFollowed(user.ID, getUserInfoDTO.UserID)

	res.StatusCode = response.Success
	app.ZeroCheck(res.User.ID)
	c.JSON(200, res)
}
