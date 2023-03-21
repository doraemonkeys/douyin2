package favorite

import (
	"strconv"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/msgQueue"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const FavoriteSuccess = "点赞成功"

type PostFavorDTO struct {
	VideoID int    `json:"video_id"`
	Token   string `json:"token"`
	// 1-点赞，2-取消点赞
	ActionType string `json:"action_type"`
}

func PostFavorHandler(c *gin.Context) {
	var postFavorDTO PostFavorDTO
	err := c.ShouldBind(&postFavorDTO)
	if err != nil {
		logrus.Error("bind postFavorDTO failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	// 获取用户
	user := c.MustGet("user").(app.User)
	// 点赞
	if postFavorDTO.ActionType != "1" && postFavorDTO.ActionType != "2" {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	action, _ := strconv.Atoi(postFavorDTO.ActionType)

	// 发送到消息队列
	mq := msgQueue.GetFavoriteMQ()
	mq.Push(msgQueue.FavoriteMSg{
		VideoID:    uint(postFavorDTO.VideoID),
		UserID:     user.ID,
		ActionType: action,
	})
	response.ResponseSuccess(c, FavoriteSuccess)
}
