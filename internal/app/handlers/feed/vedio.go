package feed

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/middleware"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const FeedVedioListLimit = 30

type FeedVideeDTO struct {
	LatestTime string   `json:"latest_time"`
	Token      string   `json:"token"`
	User       app.User `json:"user"`
}

type proxyFeedVideoList struct {
	*gin.Context
}

func NewProxyFeedVideoList(c *gin.Context) *proxyFeedVideoList {
	return &proxyFeedVideoList{Context: c}
}

// DoNoToken 未登录的视频流推送处理
func (p *proxyFeedVideoList) DoNoToken(feedRequest FeedVideeDTO) {
	var res response.VideoListResponse

	foramtedTime, err := utils.GetFormatedTimeFromUnix(feedRequest.LatestTime, models.DataBaseTimeFormat)
	if err != nil {
		response.ResponseError(p.Context, response.ErrServerInternal)
	}
	logrus.Debug("latestTime: ", feedRequest.LatestTime, "foramtedTime: ", foramtedTime)

	videoModels, err := services.GetVideoAndAuthorListFeedByLastTime(foramtedTime, FeedVedioListLimit)
	if err != nil {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrServerInternal
		p.Context.JSON(http.StatusOK, res)
		return
	}
	var dummyMap map[uint]bool = make(map[uint]bool)
	res.SetValues(videoModels, dummyMap)
	res.CommonResponse.StatusCode = response.Success
	p.Context.JSON(http.StatusOK, res)
}

// DoHasToken 如果是登录状态，则生成UserId字段
func (p *proxyFeedVideoList) DoHasToken(feedRequest FeedVideeDTO) {
	middleware.JWTMiddleWare()(p.Context)
	_, ok := p.Context.Get("user")
	if !ok {
		// middleware已经处理了错误，这里不需要处理
		return
	}

	var res response.VideoListResponse
	foramtedTime, err := utils.GetFormatedTimeFromUnix(feedRequest.LatestTime, models.DataBaseTimeFormat)
	if err != nil {
		response.ResponseError(p.Context, response.ErrServerInternal)
	}
	logrus.Debug("latestTime: ", feedRequest.LatestTime, "foramtedTime: ", foramtedTime)

	videoModels, err := services.GetVideoAndAuthorListFeedByLastTime(foramtedTime, FeedVedioListLimit)
	if err != nil {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrServerInternal
		p.Context.JSON(http.StatusOK, res)
		return
	}
	var UserIDs []uint
	for _, video := range videoModels {
		UserIDs = append(UserIDs, video.Author.ID)
	}
	FollowedMap, err := services.QueryFollowedMapByUserIDList(feedRequest.User.ID, UserIDs)
	res.SetValues(videoModels, FollowedMap)
	res.CommonResponse.StatusCode = response.Success
	p.Context.JSON(http.StatusOK, res)
}

// 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func FeedVideoListHandler(c *gin.Context) {

	proxy := NewProxyFeedVideoList(c)

	var feedRequest FeedVideeDTO
	feedRequest.LatestTime = c.DefaultQuery("latest_time", fmt.Sprint(time.Now().Unix()))
	var ok bool
	user, ok := c.Get("user")
	//未登录
	if !ok {
		proxy.DoNoToken(feedRequest)
		return
	}
	//已登录
	feedRequest.User = user.(app.User)
	logrus.Debug("feedRequest：", feedRequest)
	proxy.DoHasToken(feedRequest)
}
