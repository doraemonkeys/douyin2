package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/Doraemonkeys/douyin2/config"
	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	myjwt "github.com/Doraemonkeys/douyin2/pkg/jwt"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var JwtAuth *myjwt.CryptJWT

var once sync.Once

func InitJwt() {
	once.Do(func() {
		jwtConfig := config.GetJwtConfig()
		SigningKey, err := utils.HexStrToBytes(jwtConfig.SignKeyHex)
		if err != nil {
			logrus.Panic("初始化jwt失败, error:" + err.Error())
		}
		var cryptoer *utils.CbcAESCrypt
		cryptoer, err = utils.NewAESCrypt(jwtConfig.SecretHex)
		if err != nil {
			logrus.Panic("初始化jwt失败, error:" + err.Error())
		}
		JwtAuth = myjwt.NewJWT(SigningKey, cryptoer)
	})
}

// JWTMiddleWare 鉴权中间件，鉴权并设置user
func JWTMiddleWare(omitPaths ...string) gin.HandlerFunc {
	InitJwt()

	return func(c *gin.Context) {
		tokenStr, ok := c.GetQuery("token")
		//如果是忽略的路径，直接跳过
		for _, path := range omitPaths {
			if c.FullPath() == path && !ok {
				logrus.Info(path, "跳过鉴权")
				c.Next()
				return
			}
		}
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//验证token
		CustomClaims, err := JwtAuth.ParseToken(tokenStr)
		//如果token过期，返回错误信息
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusOK, response.CommonResponse{
				StatusCode: response.TokenExpired,
				StatusMsg:  err.Error(),
			})
			c.Abort() //阻止执行
			return
		}
		//如果token无效，返回错误信息
		if err != nil {
			logrus.Debug("token无效", err)
			c.JSON(http.StatusOK, response.CommonResponse{
				StatusCode: response.Failed,
				StatusMsg:  jwt.ErrTokenInvalidId.Error(),
			})
			c.Abort() //阻止执行
			return
		}

		var user app.User
		id, _ := strconv.ParseUint(CustomClaims.ID, 10, 64)
		user.ID = uint(id)
		user.Username = CustomClaims.Username
		user.Claims = CustomClaims
		c.Set("user", user)
		c.Next()
	}
}

func CreateToken(id uint, username string) (string, error) {
	claims := myjwt.CustomClaims{
		Username: username,
	}
	claims.ID = strconv.FormatUint(uint64(id), 10)
	return JwtAuth.CreateToken(claims)
}
