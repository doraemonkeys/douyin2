package msgQueue

import (
	"errors"
	"sync"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/Doraemonkeys/douyin2/internal/pkg/messagequeue"
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

var favoriteMQ *messagequeue.SimpleMQ[FavoriteMSg]
var favoriteMQInitOnce sync.Once

// GetFavoriteMQ
// 获取点赞消息队列
func GetFavoriteMQ() messagequeue.MQ[FavoriteMSg] {
	return favoriteMQ
}

// 点赞消息队列
func InitFavoriteMQ(msgHandler func(FavoriteMSg) error) {
	favoriteMQInitOnce.Do(func() {
		favoriteMQ = messagequeue.NewSimpleMQ(FavoriteWorkerNum, msgHandler)
	})
}

func FavoriteMsgHandler(msg FavoriteMSg) error {
	db := database.GetMysqlDB()
	like_count := models.VideoModelTable_LikeCount
	if msg.ActionType == 1 {
		// 点赞
		err := db.Model(&models.VideoModel{}).Where("id = ?", msg.VideoID).Update(like_count, db.Raw(like_count+" + ?", 1)).Error
		if err != nil {
			return err
		}
	} else if msg.ActionType == 2 {
		// 取消点赞
		err := db.Model(&models.VideoModel{}).Where("id = ?", msg.VideoID).Update(like_count, db.Raw(like_count+" - ?", 1)).Error
		if err != nil {
			return err
		}
	}
	return errors.New(ErrParam)
}
