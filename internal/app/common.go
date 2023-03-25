package app

import (
	"github.com/Doraemonkeys/douyin2/config"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type User struct {
	models.UserModel
	Token  string            `json:"token"`
	Claims *jwt.CustomClaims `json:"claims"`
}

const UserKeyName = "user"

func ZeroCheck[T comparable](v ...T) bool {
	if !config.IsDebug() {
		return false
	}
	var zero T
	for _, item := range v {
		if item == zero {
			logrus.Errorf("zero value: %v", v)
		}
	}
	return false
}
