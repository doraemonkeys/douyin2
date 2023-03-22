package database

import (
	"fmt"

	"github.com/Doraemonkeys/douyin2/config"
	"github.com/Doraemonkeys/douyin2/internal/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetMysqlDB() *gorm.DB {
	return db
}

func init() {
	connectMysql()
}

func connectMysql() {
	mysqlConf := config.GetMysqlConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.Dbname,
		mysqlConf.Timeout,
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error:" + err.Error())
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(config.GetMysqlConfig().MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.GetMysqlConfig().MaxOpenConns)
	mirateTable()
}

func mirateTable() {
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_FollowersSlice, &models.UserFollowerModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_FansSlice, &models.UserFanModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_LikesSlice, &models.UserLikeModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_CollectionsSlice, &models.UserCollectionModel{})

	db.SetupJoinTable(&models.VideoModel{}, models.VideoModelTable_LikesSlice, &models.UserLikeModel{})
	db.SetupJoinTable(&models.VideoModel{}, models.VideoModelTable_CollectionsSlice, &models.UserCollectionModel{})

	db.AutoMigrate(
		&models.UserModel{},
		&models.VideoModel{},
		&models.CommentModel{},
		&models.UserFollowerModel{},
		&models.UserFanModel{},
		&models.UserLikeModel{},
		&models.UserCollectionModel{},
	)
}
