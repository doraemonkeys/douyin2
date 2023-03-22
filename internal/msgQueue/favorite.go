package msgQueue

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/pkg/messageQueue"
	"github.com/sirupsen/logrus"
)

const (
	// 函数传入的参数错误
	ErrParam = "不合法的参数"
)

type FavoriteMSg struct {
	VideoID uint `json:"video_id"`
	UserID  uint `json:"user_id"`
	// 1-点赞，2-取消点赞
	ActionType int `json:"action_type"`
}

const FavoriteWorkerNum int = 5

var favoriteMQ *messageQueue.SimpleMQ[FavoriteMSg]
var favoriteMQInitOnce sync.Once

// GetFavoriteMQ
// 获取点赞消息队列
func GetFavoriteMQ() messageQueue.MQ[FavoriteMSg] {
	return favoriteMQ
}

// 点赞消息队列
func InitFavoriteMQ(msgHandler func(FavoriteMSg)) {
	favoriteMQInitOnce.Do(func() {
		favoriteMQ = messageQueue.NewSimpleMQ(FavoriteWorkerNum, msgHandler)
	})
}

func FavoriteMsgHandler(msg FavoriteMSg) {
	if msg.ActionType == 1 {
		// 点赞
		err := services.LikeVideo(msg.UserID, msg.VideoID)
		if err != nil {
			logrus.Error("点赞失败：", err)
		}
	} else if msg.ActionType == 2 {
		// 取消点赞
		err := services.DislikeVideo(msg.UserID, msg.VideoID)
		if err != nil {
			logrus.Error("取消点赞失败：", err)
		}
	} else {
		logrus.Error("不合法的参数：", msg)
	}
}
