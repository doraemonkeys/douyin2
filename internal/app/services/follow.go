package services

import (
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"gorm.io/gorm"
)

func FollowUser(userID, toUserID uint) error {
	// 更新db - follow表
	db := database.GetMysqlDB()
	tx := db.Begin()

	var follow models.UserFollowerModel
	follow.UserID = userID
	follow.FollowerID = toUserID
	err := tx.Create(&follow).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 更新db - user表
	var user models.UserModel
	user.ID = userID
	follower_count := models.UserModelTable_FollowerCount
	err = tx.Model(&user).Update(follower_count, gorm.Expr(follower_count+"+?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var toUser models.UserModel
	toUser.ID = toUserID
	fan_count := models.UserModelTable_FanCount
	err = tx.Model(&toUser).Update(fan_count, gorm.Expr(fan_count+"+?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 更新缓存
	userCacher := database.GetUserInfoCacher()
	userCache, exist := userCacher.Get(userID)
	if exist {
		userCache.FollowerCount++
		userCacher.Set(userID, userCache)
	}
	toUserCache, exist := userCacher.Get(toUserID)
	if exist {
		toUserCache.FanCount++
		userCacher.Set(toUserID, toUserCache)
	}
	return nil
}

func UnfollowUser(userID, toUserID uint) error {
	// 更新db - follow表
	db := database.GetMysqlDB()
	tx := db.Begin()

	var follow models.UserFollowerModel
	follow.UserID = userID
	follow.FollowerID = toUserID
	err := tx.Delete(&follow).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 更新db - user表
	var user models.UserModel
	user.ID = userID
	follower_count := models.UserModelTable_FollowerCount
	err = tx.Model(&user).Update(follower_count, gorm.Expr(follower_count+"-?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var toUser models.UserModel
	toUser.ID = toUserID
	fan_count := models.UserModelTable_FanCount
	err = tx.Model(&toUser).Update(fan_count, gorm.Expr(fan_count+"-?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 更新缓存
	userCacher := database.GetUserInfoCacher()
	userCache, exist := userCacher.Get(userID)
	if exist {
		userCache.FollowerCount--
		userCacher.Set(userID, userCache)
	}
	toUserCache, exist := userCacher.Get(toUserID)
	if exist {
		toUserCache.FanCount--
		userCacher.Set(toUserID, toUserCache)
	}
	return nil
}
