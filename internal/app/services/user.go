package services

import (
	"errors"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

/*
   Query 从MySQL中查询
   Get 从Cache中或MySQL中查询
   Query 后都会尝试存入Cache
*/

func QueryUserById(id uint) (models.UserModel, error) {
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New(response.ErrUserNotExists)
		}
		logrus.Error("query user failed, err: ", err)
		return user, errors.New(response.ErrServerInternal)
	}
	// 存入缓存
	cacher := database.GetUserInfoCacher()
	var userCacheModel models.UserCacheModel
	userCacheModel.SetValue(user)
	cacher.Set(id, userCacheModel)
	return user, nil
}

func GetUserById(id uint) (models.UserModel, error) {
	var userReturn models.UserModel
	// 从缓存中获取
	cacher := database.GetUserInfoCacher()
	user, exist := cacher.Get(id)
	if exist {
		userReturn.SetValueFromCacheModel(user)
		return userReturn, nil
	}
	// 从MySQL中获取
	var err error
	userReturn, err = QueryUserById(id)
	if err != nil {
		return userReturn, err
	}
	return userReturn, nil
}

func QueryUserByUsername(username string) (models.UserModel, error) {
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Where(models.UserModelTable_Username+" = ?", username).Find(&user).Error
	if err != nil || user.ID == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) || user.ID == 0 {
			return user, errors.New(response.ErrUserNotExists)
		}
		logrus.Error("query user failed, err: ", err)
		return user, errors.New(response.ErrServerInternal)
	}
	return user, nil
}

// QueryUserExistById 查询MySQL中是否存在某个用户。
// 不更新缓存。
func QueryUserExistById(id int) bool {
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Unscoped().Where("id = ?", id).Find(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && user.ID > 0 {
		return true
	}
	return false
}

// QueryUserExistByUsername 查询MySQL中是否存在某个用户。
// 不更新缓存。
func QueryUserExistByUsername(username string) bool {
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Unscoped().Where(models.UserModelTable_Username+" = ?", username).Find(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && user.ID > 0 {
		return true
	}
	return false
}

// QueryUserWithFollowers 用户信息与粉丝列表
// func QueryUserWithFollowersByID(id int) (models.UserModel, error) {
// 	const fieldFollower = models.UserModelTable_FollowersSlice
// 	var user models.UserModel
// 	db := database.GetMysqlDB()
// 	err := db.Debug().Preload(fieldFollower).Where("id = ?", id).Find(&user).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return user, errors.New(response.ErrUserNotExists)
// 		}
// 		logrus.Error("query user failed, err: ", err)
// 		return user, errors.New(response.ErrServerInternal)
// 	}
// 	return user, nil
// }

// 判断某个用户是否关注了另一个用户
func QueryUserFollowed(userID uint, followID uint) bool {
	var UserFollows models.UserFollowerModel
	db := database.GetMysqlDB()
	err := db.Model(&UserFollows).
		Where(models.UserFollowerModelTable_UserID+" = ? AND "+
			models.UserFollowerModelTable_FollowerID+" = ?", userID, followID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		logrus.Error("query user failed, err: ", err)
		return false
	}
	return true
}

// 判断某个用户关注了列表中的哪些用户
func QueryUserFollowedMap(userID uint, followIDList []uint) (map[uint]bool, error) {
	var UserFollows []models.UserFollowerModel
	db := database.GetMysqlDB()
	err := db.Model(&models.UserFollowerModel{}).
		Where(models.UserFollowerModelTable_UserID+" = ? AND "+
			models.UserFollowerModelTable_FollowerID+" IN (?)", userID, followIDList).Find(&UserFollows).Error

	logrus.Debug("UserFollows: ", UserFollows)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logrus.Error("query user failed, err: ", err)
		return nil, err
	}

	followedMap := make(map[uint]bool)
	for _, userFollow := range UserFollows {
		followedMap[userFollow.FollowerID] = true
	}
	return followedMap, nil
}

func QueryFavorVideoIDListByUserID(userID uint) (likeVideos []uint, err error) {
	db := database.GetMysqlDB()
	var userFavor []models.UserLikeModel
	db.Where(models.UserLikeModelTable_UserID+" = ?", userID).Find(&userFavor)
	for _, userFavor := range userFavor {
		likeVideos = append(likeVideos, userFavor.VideoID)
	}
	return likeVideos, nil
}

func QueryUserMapsByUserIDList(userIDList []uint) (userList map[uint]models.UserModel, err error) {
	userList = make(map[uint]models.UserModel)
	db := database.GetMysqlDB()
	var userListTemp []models.UserModel
	err = db.Where("id IN (?)", userIDList).Find(&userListTemp).Error
	if err != nil {
		logrus.Error("query user failed, err: ", err)
		return userList, errors.New(response.ErrServerInternal)
	}
	cacher := database.GetUserInfoCacher()
	var tempCache models.UserCacheModel
	for _, user := range userListTemp {
		userList[user.ID] = user
		tempCache.SetValue(user)
		// 更新缓存
		cacher.Set(user.ID, tempCache)
	}
	return userList, nil
}

func GetUserMapByUserIdMap[T any](userIdMap map[uint]T) (userMap map[uint]models.UserModel, err error) {
	// from cache
	userCache := database.GetUserInfoCacher()
	userMap = make(map[uint]models.UserModel)
	var temp models.UserModel
	cacheMissIDs := make([]uint, 0)
	for userId := range userIdMap {
		if user, exist := userCache.Get(userId); exist {
			temp.SetValueFromCacheModel(user)
			userMap[userId] = temp
		} else {
			cacheMissIDs = append(cacheMissIDs, userId)
		}
	}
	// from db
	if len(cacheMissIDs) > 0 {
		var userList []models.UserModel
		userList, err = QueryUserListByUserIDList(cacheMissIDs)
		if err != nil {
			logrus.Error("query user failed, err: ", err)
			return userMap, errors.New(response.ErrServerInternal)
		}
		for _, user := range userList {
			userMap[user.ID] = user
		}
	}
	return userMap, nil
}

func QueryUserListByUserIDList(userIDList []uint) (userList []models.UserModel, err error) {
	db := database.GetMysqlDB()
	app.ZeroCheck(userIDList...)
	err = db.Where("id IN (?)", userIDList).Find(&userList).Error
	if err != nil {
		logrus.Error("query user failed, err: ", err)
		return userList, errors.New(response.ErrServerInternal)
	}
	// 更新缓存
	cacher := database.GetUserInfoCacher()
	var tempCache models.UserCacheModel
	for _, user := range userList {
		tempCache.SetValue(user)
		cacher.Set(user.ID, tempCache)
	}
	return userList, nil
}

func QueryFollowedMapByUserIDList(id uint, userIDList []uint) (followedMap map[uint]bool, err error) {
	followedMap = make(map[uint]bool, len(userIDList))
	if len(userIDList) == 0 || userIDList == nil {
		return followedMap, nil
	}
	db := database.GetMysqlDB()
	var userFollows []models.UserFollowerModel
	err = db.Debug().Where(models.UserFollowerModelTable_UserID+" = ? AND "+
		models.UserFollowerModelTable_FollowerID+" IN (?)", id, userIDList).Find(&userFollows).Error
	if err != nil {
		logrus.Error("query user failed, err: ", err)
		return followedMap, errors.New(response.ErrServerInternal)
	}
	for _, userFollow := range userFollows {
		followedMap[userFollow.FollowerID] = true
	}
	return followedMap, nil
}

func QueryFollowedMapByUserIDMap[T any](id uint, userIDMap map[uint]T) (followedMap map[uint]bool, err error) {
	followedMap = make(map[uint]bool, len(userIDMap))
	db := database.GetMysqlDB()
	var userFollows []models.UserFollowerModel
	var ids []uint
	for id := range userIDMap {
		ids = append(ids, id)
	}
	err = db.Debug().Where(models.UserFollowerModelTable_UserID+" = ? AND "+
		models.UserFollowerModelTable_FollowerID+" IN (?)", id, ids).Find(&userFollows).Error
	if err != nil {
		logrus.Error("query user failed, err: ", err)
		return followedMap, errors.New(response.ErrServerInternal)
	}
	for _, userFollow := range userFollows {
		followedMap[userFollow.FollowerID] = true
	}
	return followedMap, nil
}
