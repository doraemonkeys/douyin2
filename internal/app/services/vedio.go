package services

import (
	"errors"
	"math"
	"strings"

	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Doraemonkeys/douyin2/internal/database"
)

/*
   Query 从MySQL中查询
   Get 从Cache中或MySQL中查询
   Query 后都会尝试存入Cache
*/

const (
	// 函数传入的参数错误
	ErrParam = "不合法的参数"
	// 重复插入
	ErrDuplicate = "重复插入"
	// 删除不存在的记录
	ErrDeleteNotExists = "删除不存在的记录"
)

const (
	// 重复插入
	MysqlDuplicatePrefix = "Error 1062"
)

// QueryVideoAndUserListByLastTime
// 查询MySQL，返回按投稿时间倒序的视频列表，视频数由服务端控制，limit为最大值，
// foramtedTime 为格式化后的时间，可能为：2023-03-20 04:16:17.648。
// 查询的数据会尝试存入Cache。
func QueryVideoAndUserListByLastTime(foramtedTime string, limit int) ([]models.VideoModel, error) {
	if limit <= 0 {
		return nil, errors.New(ErrParam)
	}
	var videos []models.VideoModel
	db := database.GetMysqlDB()
	created_at := models.VideoModelTable_CreatedAt
	err := db.Debug().Where(created_at+" <= ?", foramtedTime).Order(created_at + " desc").Limit(limit).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	logrus.Trace("QueryVideoList：", videos)

	// videoInfo 存入Cache
	videoCacher := database.GetVideoInfoCacher()
	for _, v := range videos {
		var temp models.VideoCacheModel
		temp.SetValue(v)
		videoCacher.Set(v.ID, temp)
	}

	var ids []uint
	for _, val := range videos {
		ids = append(ids, val.AuthorID)
	}

	// Query
	var Users []models.UserModel
	err = db.Debug().Where("id in (?)", ids).Find(&Users).Error
	if err != nil {
		logrus.Error(err)
		return nil, errors.New(response.ErrServerInternal)
	}
	logrus.Debug("QueryVideoUserList：", Users)

	// 将Users信息添加到videos中
	var UsersMap = make(map[uint]models.UserModel)
	for _, v := range Users {
		UsersMap[v.ID] = v
	}
	for _, v := range videos {
		v.Author = UsersMap[v.AuthorID]
	}
	// userInfo 存入Cache
	userCacher := database.GetUserInfoCacher()
	for _, v := range Users {
		var temp models.UserCacheModel
		temp.SetValue(v)
		userCacher.Set(v.ID, temp)
	}

	return videos, nil
}

// GetVideoAndAuthorListFeedByLastTime 返回的VideoModel切片中包含Author信息
// 查询MySQL和Cache，返回视频列表(热门limit*2/3 + 新发布limit*1/3)，视频数由服务端控制，limit为最大值。
// 注意不能使用返回的结构体中的指针或引用类型，可能为nil。
func GetVideoAndAuthorListFeedByLastTime(foramtedTime string, limit int) ([]models.VideoModel, error) {
	// 从Cache中获取符合要求的视频列表(热门limit*(2/3))
	var hotRate = 2.0 / 3.0
	// 热门视频数，向上取整
	var hotLimit = math.Ceil(float64(limit) * hotRate)
	cacher := database.GetVideoInfoCacher()
	videosCache, err := cacher.GetRandomMulti(int(hotLimit))
	if err != nil {
		return nil, err
	}
	// 向videos1中添加作者信息
	// 1. from cache
	var videos1 []models.VideoModel
	userCacher := database.GetUserInfoCacher()
	var cacheMiss = make([]uint, 0)
	for _, v := range videosCache {
		var temp models.VideoModel
		temp.SetValueFromCacheModel(v)
		user, exist := userCacher.Get(v.AuthorID)
		if exist {
			temp.Author.SetValueFromCacheModel(user)
		} else {
			cacheMiss = append(cacheMiss, v.AuthorID)
		}
		videos1 = append(videos1, temp)
	}
	// 2. from mysql 处理cache miss
	if len(cacheMiss) > 0 {
		Users, err := QueryUserListByUserIDList(cacheMiss)
		if err != nil {
			return nil, err
		}
		// 将Users信息添加到videos1中
		var UsersMap = make(map[uint]models.UserModel)
		for _, v := range Users {
			UsersMap[v.ID] = v
		}
		for _, v := range videos1 {
			if user, exist := UsersMap[v.AuthorID]; exist {
				v.Author = user
			}
		}
	}
	newLimit := float64(limit) - float64(len(videos1))
	var videos2 []models.VideoModel
	//var UsersMap = make(map[uint]models.UserModel)
	// 从MySQL中获取符合要求的视频列表(新发布limit*(1/3))
	if newLimit > 0 {
		videos2, err = QueryVideoAndUserListByLastTime(foramtedTime, int(newLimit))
		if err != nil {
			return nil, err
		}
	}
	// 合并两个视频列表
	Videos := append(videos1, videos2...)
	return Videos, nil
}

func CreateVedio(video *models.VideoModel) error {
	db := database.GetMysqlDB()
	err := db.Create(video).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryPublishListByAuthorID(userID uint) ([]models.VideoModel, error) {
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

// 点赞视频
func LikeVideo(userID, videoID uint) error {
	// 1. 更新userLike缓存
	err := UpdateUserLikeCache(userID, videoID, true)
	if err != nil {
		return err
	}

	db := database.GetMysqlDB()
	// 2. 更新userLike表
	var userLike models.UserLikeModel
	userLike.UserID = userID
	userLike.VideoID = videoID
	err = db.Create(&userLike).Error
	if err != nil {
		if strings.HasPrefix(err.Error(), MysqlDuplicatePrefix) {
			return errors.New(ErrDuplicate)
		}
		return err
	}
	// 3. 更新video表
	like_count := models.VideoModelTable_LikeCount
	err = db.Model(&models.VideoModel{}).Where("id = ?", videoID).
		Update(like_count, gorm.Expr(like_count+" + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserLikeCache 更新用户点赞缓存。
// isLike为true表示点赞，false表示取消点赞。
// 缓存未命中时，直接返回nil。
// 取消不存在的点赞时，返回ErrDeleteNotExists。
func UpdateUserLikeCache(userID, videoID uint, isLike bool) error {
	cacher := database.GetUserFavoriteCacher()
	// 1. 查询userLike缓存,如果不存在则直接返回
	likes, exist := cacher.Get(userID)
	if !exist {
		return nil
	}
	// 2. 更新userLike缓存和video缓存(如果存在)
	if isLike {
		// 点赞
		// 从cache中查询该视频的信息
		var NewLikeCache = models.UserLike_VideoAndAuthor{}
		videoCacher := database.GetVideoInfoCacher()
		videoCache, videoCachExist := videoCacher.Get(videoID)
		if !videoCachExist {
			// 从MySQL中查询该视频的信息
			video, err := QueryVideoInfoByID(videoID)
			if err != nil {
				return err
			}
			NewLikeCache.AuthorID = video.ID
		} else {
			NewLikeCache.AuthorID = videoCache.AuthorID
			videoCache.LikeCount++
			// 更新video缓存
			videoCacher.Set(videoID, videoCache)
		}
		NewLikeCache.IsFollowed = QueryUserFollowed(userID, videoCache.AuthorID)
		likes.VideoIDMap.Store(videoID, NewLikeCache)
		return nil
	}
	// 取消点赞
	_, likeExist := likes.VideoIDMap.LoadAndDelete(videoID)
	if !likeExist {
		return errors.New(ErrDeleteNotExists)
	}
	// 从cache中查询该视频的信息并更新
	videoCacher := database.GetVideoInfoCacher()
	videoCache, videoCachExist := videoCacher.Get(videoID)
	if !videoCachExist {
		return nil
	}
	videoCache.LikeCount--
	videoCacher.Set(videoID, videoCache)
	return nil
}

func QueryVideoInfoByID(videoID uint) (models.VideoModel, error) {
	var video models.VideoModel
	db := database.GetMysqlDB()
	err := db.Where("id = ?", videoID).Take(&video).Error
	if err != nil {
		return video, err
	}
	return video, nil
}

// 取消点赞视频
func DislikeVideo(userID, videoID uint) error {

	// 1. 更新userLike缓存
	err := UpdateUserLikeCache(userID, videoID, false)
	if err != nil {
		return err
	}

	db := database.GetMysqlDB()
	// 1. 查询userLike表 是否已经点赞
	var userLike models.UserLikeModel
	user_id := models.UserLikeModelTable_UserID
	video_id := models.UserLikeModelTable_VideoID
	err = db.Where(user_id+" = ? and "+video_id+" = ?", userID, videoID).Take(&userLike).Error
	if err == gorm.ErrRecordNotFound || userLike.UserID == 0 {
		return errors.New(ErrDeleteNotExists)
	}
	// 2. 更新userLike表
	err = db.Delete(&userLike).Error
	if err != nil {
		return err
	}
	// 3. 更新video表
	like_count := models.VideoModelTable_LikeCount
	err = db.Model(&models.VideoModel{}).Where("id = ?", videoID).
		Update(like_count, gorm.Expr(like_count+" - ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryVideoListByVideoIDList(videoIDList []uint) ([]models.VideoModel, error) {
	var videos []models.VideoModel
	db := database.GetMysqlDB()
	err := db.Where("id in (?)", videoIDList).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	// 将从MySQL中获取的视频信息添加到cache中
	return nil, nil
}

func GetVideoListAndAuthorByVideoIDList(videoIDList []uint) ([]models.VideoModel, error) {
	var videos []models.VideoModel
	// 从cache中获取视频信息
	videoCacher := database.GetVideoInfoCacher()
	VideoCache := videoCacher.GetMulti(videoIDList)
	for _, v := range VideoCache {
		var temp models.VideoModel
		temp.SetValueFromCacheModel(v)
		videos = append(videos, temp)
	}
	// cache中没有的视频ID
	var videoIDCacheMiss []uint
	for _, v := range videoIDList {
		if VideoCache[v].ID == 0 {
			videoIDCacheMiss = append(videoIDCacheMiss, v)
		}
	}
	// 从MySQL中获取cache未命中的视频信息
	if len(videoIDCacheMiss) > 0 {
		MysqlVideos, err := QueryVideoListByVideoIDList(videoIDCacheMiss)
		if err != nil {
			return nil, err
		}
		videos = append(videos, MysqlVideos...)
	}

	// 从Cache中获取作者信息
	var authorIDList []uint
	for _, v := range videos {
		authorIDList = append(authorIDList, v.AuthorID)
	}
	authorCacher := database.GetUserInfoCacher()
	authorCache := authorCacher.GetMulti(authorIDList)
	for _, v := range videos {
		if authorCache[v.AuthorID].ID != 0 {
			v.Author.SetValueFromCacheModel(authorCache[v.AuthorID])
		}
	}
	// cache中没有的作者ID
	var authorIDCacheMiss []uint
	for _, v := range authorIDList {
		if authorCache[v].ID == 0 {
			authorIDCacheMiss = append(authorIDCacheMiss, v)
		}
	}
	// 从MySQL中获取cache未命中的作者信息
	if len(authorIDCacheMiss) > 0 {
		UsersMaps, err := QueryUserMapsByUserIDList(authorIDCacheMiss)
		if err != nil {
			return nil, err
		}
		for _, v := range videos {
			if v.Author.ID == 0 {
				v.Author = UsersMaps[v.AuthorID]
			}
		}
	}
	return videos, nil
}
