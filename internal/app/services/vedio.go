package services

import (
	"errors"
	"math"

	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Doraemonkeys/douyin2/internal/database"
)

/*
   Query 从MySQL中查询
   Get 从Cache中或MySQL中查询
*/

const (
	// 函数传入的参数错误
	ErrParam = "不合法的参数"
)

// QueryVideoAndUserListByLastTime
// 查询MySQL，返回按投稿时间倒序的视频列表，视频数由服务端控制，limit为最大值，
// foramtedTime 为格式化后的时间，可能为：2023-03-20 04:16:17.648。
func QueryVideoAndUserListByLastTime(foramtedTime string, limit int) ([]models.VideoModel, map[uint]models.UserModel, error) {
	if limit <= 0 {
		return nil, nil, errors.New(ErrParam)
	}
	var videos []models.VideoModel
	db := database.GetMysqlDB()
	created_at := models.VideoModelTable_CreatedAt
	err := db.Debug().Where(created_at+" <= ?", foramtedTime).Order(created_at + " desc").Limit(limit).Find(&videos).Error
	if err != nil {
		return nil, nil, err
	}
	logrus.Trace("QueryVideoList：", videos)

	var UsersMap = make(map[uint]models.UserModel)
	for _, v := range videos {
		UsersMap[v.AuthorID] = models.UserModel{}
	}
	var ids []uint
	for k := range UsersMap {
		ids = append(ids, k)
	}

	// Query
	var Users []models.UserModel
	err = db.Debug().Where("id in (?)", ids).Find(&Users).Error
	if err != nil {
		logrus.Error(err)
		return nil, nil, errors.New(response.ErrServerInternal)
	}
	logrus.Debug("QueryVideoUserList：", Users)

	for _, v := range Users {
		UsersMap[v.ID] = v
	}
	return videos, UsersMap, nil
}

// GetVideoAndUserListFeedByLastTime
// 查询MySQL和Cache，返回视频列表(热门limit*2/3 + 新发布limit*1/3)，视频数由服务端控制，limit为最大值。
// 注意不能使用返回的结构体中的指针或引用类型，可能为nil。
func GetVideoAndUserListFeedByLastTime(foramtedTime string, limit int) ([]models.VideoModel, map[uint]models.UserModel, error) {
	// 从Cache中获取符合要求的视频列表(热门limit*(2/3))
	var hotRate = 2.0 / 3.0
	// 热门视频数，向上取整
	var hotLimit = math.Ceil(float64(limit) * hotRate)
	cacher := database.GetVideoInfoCacher()
	videos1, err := cacher.GetRandomMulti(int(hotLimit))
	if err != nil {
		return nil, nil, err
	}
	newLimit := float64(limit) - float64(len(videos1))
	var videos2 []models.VideoModel
	var UsersMap = make(map[uint]models.UserModel)
	// 从MySQL中获取符合要求的视频列表(新发布limit*(1/3))
	if newLimit > 0 {
		videos2, UsersMap, err = QueryVideoAndUserListByLastTime(foramtedTime, int(newLimit))
		if err != nil {
			return nil, nil, err
		}
	}
	// 合并两个视频列表
	var videos []models.VideoModel
	for _, v := range videos1 {
		var temp models.VideoModel
		temp.SetValueFromCacheModel(v)
		videos = append(videos, temp)
	}
	videos = append(videos, videos2...)
	// 向UsersMap中添加hot视频的作者信息
	for _, v := range videos1 {
		var temp models.UserModel
		temp.SetValueFromCacheModel(v.Author)
		UsersMap[v.AuthorID] = temp
	}
	return videos, UsersMap, nil
}

func CreateVedio(video *models.VideoModel) error {
	db := database.GetMysqlDB()
	err := db.Create(video).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryVideoListByUserID(userID uint) ([]models.VideoModel, error) {
	var videos []models.VideoModel
	db := database.GetMysqlDB()
	author_id := models.VideoModelTable_AuthorID
	created_at := models.VideoModelTable_CreatedAt
	err := db.Where(author_id+" = ?", userID).Order(created_at + " desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	logrus.Debug("QueryVideoListByUserID：", userID)
	for _, v := range videos {
		logrus.Debug(v)
	}
	return videos, nil
}

// 返回用户喜欢的视频列表和用户自己的信息
func QueryLikeVideoListByUserID(userId uint) ([]models.VideoModel, models.UserModel, error) {
	var user models.UserModel
	db := database.GetMysqlDB()
	Likes := models.UserModelTable_LikesSlice
	err := db.Preload(Likes).Where("id = ?", userId).Take(&user).Error
	if err == gorm.ErrRecordNotFound || user.ID == 0 {
		return nil, user, errors.New(response.ErrUserNotExists)
	}
	if err != nil {
		return nil, user, err
	}
	return user.Likes, user, nil
}
