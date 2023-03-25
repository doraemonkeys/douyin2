package services

import (
	"errors"

	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/sirupsen/logrus"
)

func CommentVideo(videoId uint, commenterID uint, commentText string) error {
	if videoId == 0 || commenterID == 0 {
		logrus.Error("comment video failed, videoId or commenterID is 0, videoId: ", videoId, " commenterID: ", commenterID)
	}
	var comment models.CommentModel
	comment.VideoID = videoId
	comment.Content = commentText
	comment.UserID = commenterID
	comment.Commenter.ID = commenterID
	err := database.GetMysqlDB().Create(&comment).Error
	if err != nil {
		return err
	}
	AddVideoCommentToCache(videoId, comment)
	return nil
}

func AddVideoCommentToCache(videoId uint, comment models.CommentModel) {
	cacher := database.GetVideoCommentCacher()
	commentChche, exist := cacher.Get(videoId)
	if exist {
		commentChche.MapLock.Lock()
		commentChche.CacheMap[comment.ID] = comment
		commentChche.MapLock.Unlock()
		return
	}
}

func DeleteVideoCommentFromCache(videoId uint, commentId uint, commenterID uint) {
	cacher := database.GetVideoCommentCacher()
	commentChche, exist := cacher.Get(videoId)
	if exist {
		commentChche.MapLock.RLock()
		comment, exist := commentChche.CacheMap[commentId]
		commentChche.MapLock.RUnlock()
		if !exist || comment.UserID != commenterID {
			logrus.Error("delete comment failed, comment not exists or not owner ", commentId, commenterID, " exist: ", exist, " comment: ", comment, "")
			return
		}
		commentChche.MapLock.Lock()
		delete(commentChche.CacheMap, commentId)
		commentChche.MapLock.Unlock()
		return
	}
}

func DeleteComment(commentId uint, commenterID uint) error {
	var comment models.CommentModel
	err := database.GetMysqlDB().Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		return err
	}
	if comment.ID == 0 {
		return errors.New(ErrDeleteNotExists)
	}
	if comment.UserID != commenterID {
		return errors.New(ErrDeleteNotOwner)
	}
	err = database.GetMysqlDB().Delete(&comment).Error
	if err != nil {
		return err
	}
	DeleteVideoCommentFromCache(comment.VideoID, commentId, commenterID)
	return nil
}

func QueryCommentListWithCommenterByVideoID(videoId uint) ([]models.CommentModel, error) {
	db := database.GetMysqlDB()
	var commentList []models.CommentModel
	Commenter := models.CommentModelPreload_Commenter
	video_id := models.CommentModelTable_VideoID
	err := db.Debug().Preload(Commenter).Where(video_id+" = ?", videoId).Find(&commentList).Error
	if err != nil {
		logrus.Error("query comment list failed, err: ", err)
		return commentList, err
	}
	logrus.Info("query comment list success, comment list: ", commentList)
	// 存入缓存
	cacher := database.GetVideoCommentCacher()
	var commentCache = models.NewCommentCacheModel()
	for _, comment := range commentList {
		commentCache.CacheMap[comment.ID] = comment
	}
	cacher.Set(videoId, commentCache)
	return commentList, nil
}
