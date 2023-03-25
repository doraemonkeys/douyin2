package msgQueue

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/pkg/messageQueue"
	"github.com/sirupsen/logrus"
)

type CommentMsg struct {
	VideoID    uint   `json:"video_id"`
	ActionType string `json:"action_type"`
	//用户填写的评论内容，在action_type=1的时候使用
	CommentText string `json:"comment_text"`
	//要删除的评论id，在action_type=2的时候使用
	CommentId uint `json:"comment_id"`
	// 评论者或者删除者的id
	CommenterID uint `json:"commenter_id"`
}

const (
	ActionTypeComment = "1"
	ActionTypeDelete  = "2"
)

var commentMQ *messageQueue.SimpleMQ[CommentMsg]
var commentMQInitOnce sync.Once

func GetCommentMQ() messageQueue.MQ[CommentMsg] {
	return commentMQ
}

const CommentWorkerNum int = 10

func InitCommentMQ() {
	commentMQInitOnce.Do(func() {
		commentMQ = messageQueue.NewSimpleMQ(CommentWorkerNum, CommentMsgHandler)
	})
}

func CommentMsgHandler(msg CommentMsg) {
	if msg.ActionType == ActionTypeComment {
		// 发表评论
		//logrus.Debug("发表评论：", "video_id:", msg.VideoID, "commenter_id:", msg.CommenterID, "comment_text:", msg.CommentText)
		err := services.CommentVideo(msg.VideoID, msg.CommenterID, msg.CommentText)
		if err != nil {
			logrus.Error("发表评论失败：", err)
		}
	} else if msg.ActionType == ActionTypeDelete {
		// 删除评论
		//logrus.Debug("删除评论：", "comment_id:", msg.CommentId, "commenter_id:", msg.CommenterID)
		err := services.DeleteComment(msg.CommentId, msg.CommenterID)
		if err != nil {
			logrus.Error("删除评论失败：", err)
		}
	} else {
		logrus.Error("不合法的参数：", msg)
	}
}
