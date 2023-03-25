package services

import (
	"errors"

	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

/*
   Query 从MySQL中查询
   Get 从Cache中或MySQL中查询
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
	return user, nil
}

func GetUserById(id uint) (models.UserModel, error) {
	var userReturn models.UserModel
	// 从缓存中获取
	cacher := database.GetUserInfoCacher()
	user, exist := cacher.Get(id)
	if !exist {
		var err error
		userReturn, err = QueryUserById(id)
		if err != nil {
			return userReturn, err
		}
		user.SetValue(userReturn)
		cacher.Set(id, user)
	}
	userReturn.SetValueFromCacheModel(user)
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

func QueryUserExistById(id int) bool {
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Unscoped().Where("id = ?", id).Find(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && user.ID > 0 {
		return true
	}
	return false
}

func QueryUserExistByUsername(username string) bool {
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Debug().Unscoped().Where(models.UserModelTable_Username+" = ?", username).Find(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && user.ID > 0 {
		return true
	}
	return false
}

// QueryUserWithFollowers 用户信息与粉丝列表
func QueryUserWithFollowersByID(id int) (models.UserModel, error) {
	const fieldFollower = models.UserModelTable_FollowersSlice
	var user models.UserModel
	db := database.GetMysqlDB()
	err := db.Debug().Preload(fieldFollower).Where("id = ?", id).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New(response.ErrUserNotExists)
		}
		logrus.Error("query user failed, err: ", err)
		return user, errors.New(response.ErrServerInternal)
	}
	return user, nil
}

// 判断某个用户是否关注了另一个用户
func QueryUserFollowed(userID uint, followID uint) bool {
	var UserFollows models.UserFollowerModel
	db := database.GetMysqlDB()
	err := db.Debug().
		Model(&UserFollows).
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
	logrus.Debug("看看这里的sql语句对不对")
	err := db.Debug().
		Model(&UserFollows).
		Where(models.UserFollowerModelTable_UserID+" = ? AND "+
			models.UserFollowerModelTable_FollowerID+" IN (?)", userID, followIDList).Error

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

func QueryFavorVideoListIDByUserID(userID uint) (likeVideos []uint, err error) {
	return nil, nil
}

func QueryUserListByUserIDList(userIDList []uint) (userList []models.UserModel, err error) {
	return nil, nil
}

func QueryUserMapsByUserIDList(userIDList []uint) (userList map[uint]models.UserModel, err error) {
	return nil, nil
}

func QueryFollowedMapByUserIDList(id uint, userIDList []uint) (followedMap map[uint]bool, err error) {
	return nil, nil
}

func QueryFollowedMapByUserIDMap[T any](id uint, userIDMap map[uint]T) (followedMap map[uint]bool, err error) {
	return nil, nil
}
