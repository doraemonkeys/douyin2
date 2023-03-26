package msgQueue

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/pkg/messageQueue"
	"github.com/sirupsen/logrus"
)

const (
	ActionType_Follow   = "1"
	ActionType_Unfollow = "2"
)

const followWorkerNum int = 10

var followMQ *messageQueue.SimpleMQ[FollowMsg]
var followMQInitOnce sync.Once

type FollowMsg struct {
	//对方用户id
	ToUserID uint `json:"to_user_id"`
	//1-关注，2-取消关注
	ActionType string `json:"action_type"`
	UserID     uint   `json:"user_id"`
}

func GetFollowMQ() messageQueue.MQ[FollowMsg] {
	return followMQ
}

func InitFollowMQ() {
	followMQInitOnce.Do(func() {
		followMQ = messageQueue.NewSimpleMQ(followWorkerNum, FollowMsgHandler)
	})
}

func FollowMsgHandler(msg FollowMsg) {
	if msg.ActionType == ActionType_Follow {
		// 关注
		err := services.FollowUser(msg.UserID, msg.ToUserID)
		if err != nil {
			logrus.Error("关注失败：", err)
		}
	} else if msg.ActionType == ActionType_Unfollow {
		// 取消关注
		err := services.UnfollowUser(msg.UserID, msg.ToUserID)
		if err != nil {
			logrus.Error("取消关注失败：", err)
		}
	} else {
		logrus.Error("不合法的参数：", msg)
	}
}
