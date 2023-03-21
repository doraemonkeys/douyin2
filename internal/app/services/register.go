package services

import (
	"errors"
	"strings"

	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/sirupsen/logrus"
)

// CheckPassword 检查密码是否符合要求
func CheckPassword(password string) bool {
	if len(password) < 6 || len(password) > 32 {
		return false
	}
	return true
}

// CheckUsername 检查用户名是否符合要求
func CheckUsername(username string) bool {
	return !strings.Contains(username, " ") && len(username) >= 2 && len(username) <= 32
}

func CreateUser(username string, rawPassword string) (response.RegisterResponse, error) {
	var res response.RegisterResponse
	if !CheckUsername(username) {
		return res, errors.New(response.ErrInvalidUsername)
	}
	if !CheckPassword(rawPassword) {
		return res, errors.New(response.ErrInvalidPassword)
	}
	if QueryUserExistByUsername(username) {
		return res, errors.New(response.ErrUserExists)
	}
	pwdHash := utils.BcryptHash(rawPassword)
	user := models.UserModel{
		Username: username,
		Password: pwdHash,
	}
	id, err := createUser(user)
	if err != nil {
		logrus.Error("create user failed, err: ", err)
		return res, errors.New(response.ErrServerInternal)
	}
	res.CommonResponse.StatusCode = response.Success
	res.UserID = int(id)
	return res, nil
}

func createUser(user models.UserModel) (uint, error) {
	db := database.GetMysqlDB()
	err := db.Create(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
