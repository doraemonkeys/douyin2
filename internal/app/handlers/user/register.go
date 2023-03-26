package user

import (
	"net/http"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RegisterUserDTO 注册用户的请求参数
type RegisterUserDTO struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

const (
	RegisterUserDTO_Username = "username"
	RegisterUserDTO_Password = "password"
)

func UserRegisterHandler(c *gin.Context) {
	var registerRequest RegisterUserDTO
	var res response.RegisterResponse
	// if err := c.ShouldBindJSON(&registerRequest); err != nil {
	// 	logrus.Debug("UserRegisterHandler error: ", err)
	// 	res.CommonResponse.StatusCode = response.Failed
	// 	res.CommonResponse.StatusMsg = response.ErrInvalidParams
	// 	c.JSON(http.StatusOK, res)
	// 	return
	// }
	var ok1, ok2 bool
	registerRequest.Username, ok1 = c.GetQuery(RegisterUserDTO_Username)
	registerRequest.Password, ok2 = c.GetQuery(RegisterUserDTO_Password)
	if !ok1 || !ok2 {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrInvalidParams
		c.JSON(http.StatusOK, res)
		return
	}

	res, err := services.CreateUser(registerRequest.Username, registerRequest.Password)
	if err != nil {
		res.CommonResponse.StatusCode = response.Failed
		switch err.Error() {
		case response.ErrUserExists:
			res.CommonResponse.StatusMsg = response.ErrUserExists
		case response.ErrInvalidPassword:
			res.CommonResponse.StatusMsg = response.ErrInvalidPassword
		case response.ErrInvalidUsername:
			res.CommonResponse.StatusMsg = response.ErrInvalidUsername
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
