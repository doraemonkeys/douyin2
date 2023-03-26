package server

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Doraemonkeys/douyin2/config"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/comment"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/favorite"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/feed"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/follow"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/publish"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/user"
	"github.com/Doraemonkeys/douyin2/internal/app/middleware"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/gin-gonic/gin"
)

type DouyinServer struct {
	Router       *gin.Engine
	PanicHandler gin.HandlerFunc
	// Service
	// ...
}

func NewDouyinServer() *DouyinServer {
	router := initDouyinRouter()
	return &DouyinServer{Router: router}
}

func (s *DouyinServer) Run(addr string) error {
	return s.Router.Run(addr)
}

func initPanicLogWriter() io.Writer {
	panicLogPath := filepath.Join(config.GetLogConfig().Path, config.GetLogConfig().PanicLogName)
	return utils.GetNewLazyFileWriter(panicLogPath)
}

// middleware -> handler -> service -> database
//
// middleware -> handler -> service -> message queue -> database
//
// middleware -> handler -> service -> cache -> database
func initDouyinRouter() *gin.Engine {
	router := gin.New()
	writer := io.MultiWriter(initPanicLogWriter(), os.Stdout)
	router.Use(gin.Logger(), gin.RecoveryWithWriter(writer))

	router.Static(config.GetVedioConfig().UrlPrefix, config.GetVedioConfig().BasePath)

	baseGroup := router.Group("/douyin")

	// basic api
	baseGroup.GET("/feed", middleware.JWTMiddleWare("/douyin/feed"), feed.FeedVideoListHandler)
	baseGroup.POST("/user/register/", user.UserRegisterHandler)
	baseGroup.POST("/user/login/", middleware.UserLoginHandler)
	baseGroup.GET("/user/", middleware.JWTMiddleWare(), user.GetUserInfoHandler)
	baseGroup.POST("/publish/action/", middleware.JWTMiddleWare(), publish.PublishVedioHandler)
	baseGroup.GET("/publish/list/", middleware.JWTMiddleWare(), publish.QueryPublishListHandler)

	//extend 1
	baseGroup.POST("/favorite/action/", middleware.JWTMiddleWare(), favorite.PostFavorHandler)
	baseGroup.GET("/favorite/list/", middleware.JWTMiddleWare(), favorite.QueryFavorVideoListHandler)
	baseGroup.POST("/comment/action/", middleware.JWTMiddleWare(), comment.PostCommentHandler)
	baseGroup.GET("/comment/list/", middleware.JWTMiddleWare(), comment.QueryCommentListHandler)

	//extend 2
	baseGroup.POST("/relation/action/", middleware.JWTMiddleWare(), follow.PostFollowActionHandler)
	// baseGroup.GET("/relation/follow/list/", middleware.NoAuthToGetUserId(), user_info.QueryFollowListHandler)
	// baseGroup.GET("/relation/follower/list/", middleware.NoAuthToGetUserId(), user_info.QueryFollowerHandler)

	return router
}
