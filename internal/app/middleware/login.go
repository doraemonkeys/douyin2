package middleware

import (
	"errors"
	"net/http"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoginRequestDto 登录请求
type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginHandler 用户登录
func UserLoginHandler(c *gin.Context) {
	var loginRequest LoginRequestDto
	var res response.LoginResponse
	// if err := c.ShouldBindJSON(&loginRequest); err != nil {
	// 	res.CommonResponse.StatusCode = response.Failed
	// 	res.CommonResponse.StatusMsg = response.ErrInvalidParams
	// 	c.JSON(http.StatusOK, res)
	// 	return
	// }
	var ok1, ok2 bool
	loginRequest.Username, ok1 = c.GetQuery("username")
	loginRequest.Password, ok2 = c.GetQuery("password")
	if !ok1 || !ok2 {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrInvalidParams
		c.JSON(http.StatusOK, res)
		return
	}
	res, err := Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		res.CommonResponse.StatusCode = response.Failed
		switch err.Error() {
		case response.ErrUserNotExists:
			res.CommonResponse.StatusMsg = response.ErrUserNotExists
		case response.ErrUserPassword:
			res.CommonResponse.StatusMsg = response.ErrUserPassword
		default:
			res.CommonResponse.StatusMsg = response.ErrServerInternal
		}
		c.JSON(http.StatusOK, res)
		return
	}
	res.CommonResponse.StatusCode = response.Success
	app.ZeroCheck(res.UserID)
	c.JSON(http.StatusOK, res)
}

func Login(username string, password string) (response.LoginResponse, error) {
	var res response.LoginResponse
	user, err := services.QueryUserByUsername(username)
	logrus.Debug("login user: ", user)
	if err != nil {
		if err.Error() == response.ErrUserNotExists {
			return res, errors.New(response.ErrUserNotExists)
		}
		logrus.Error("QueryUserByUsername error: ", err)
		return res, errors.New(response.ErrServerInternal)
	}
	if !utils.BcryptMatch(user.Password, password) {
		return res, errors.New(response.ErrUserPassword)
	}
	token, err := CreateToken(user.ID, user.Username)
	if err != nil {
		logrus.Error("CreateToken error: ", err)
		return res, errors.New(response.ErrServerInternal)
	}
	res.Token = token
	res.UserID = int(user.ID)
	return res, nil
}
