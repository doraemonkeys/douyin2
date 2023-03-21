package app

import (
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/pkg/jwt"
)

type User struct {
	models.UserModel
	Token  string            `json:"token"`
	Claims *jwt.CustomClaims `json:"claims"`
}

func VedioFilePathToUrl(path string) string {
	return "http://localhost:8080" + path
}

func CoverFilePathToUrl(path string) string {
	return "http://localhost:8080" + path
}

func VedioUrlToFilePath(url string) string {
	return url
}

func CoverUrlToFilePath(url string) string {
	return url
}
