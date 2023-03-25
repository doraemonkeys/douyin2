package feed

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/middleware"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const FeedVedioListLimit = 30

type FeedVideeDTO struct {
	// 最新的视频时间戳，毫秒
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

	lastTime, err := strconv.ParseInt(feedRequest.LatestTime, 10, 64)
	if err != nil {
		response.ResponseError(p.Context, response.ErrInvalidParams)
	}

	videoModels, err := services.GetVideoAndAuthorListFeedByLastTime(lastTime, FeedVedioListLimit)
	if err != nil && err.Error() == services.ErrDBEmpty {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrDBEmpty
		p.Context.JSON(http.StatusOK, res)
		return
	}
	if err != nil {
		logrus.Error("get video list failed, err: ", err)
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrServerInternal
		p.Context.JSON(http.StatusOK, res)
		return
	}
	//debug
	for _, video := range videoModels {
		app.ZeroCheck(video.Author.ID)
	}
	var dummyMap map[uint]bool = make(map[uint]bool)
	res.SetValues(videoModels, dummyMap, dummyMap)
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
	lastTime, err := strconv.ParseInt(feedRequest.LatestTime, 10, 64)
	if err != nil {
		response.ResponseError(p.Context, response.ErrInvalidParams)
	}
	videoModels, err := services.GetVideoAndAuthorListFeedByLastTime(lastTime, FeedVedioListLimit)
	if err != nil && err.Error() == services.ErrDBEmpty {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrDBEmpty
		p.Context.JSON(http.StatusOK, res)
		return
	}
	if err != nil {
		logrus.Error("get video list failed, err: ", err)
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrServerInternal
		p.Context.JSON(http.StatusOK, res)
		return
	}
	var UserIDs []uint
	for _, video := range videoModels {
		UserIDs = append(UserIDs, video.Author.ID)
	}
	logrus.Debug("UserIDs: ", UserIDs)
	FollowedMap, err := services.QueryFollowedMapByUserIDList(feedRequest.User.ID, UserIDs)
	if err != nil {
		logrus.Error("get followed map failed, err: ", err)
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrServerInternal
		p.Context.JSON(http.StatusOK, res)
		return
	}
	feedVIDs := make([]uint, 0, len(videoModels))
	for _, video := range videoModels {
		feedVIDs = append(feedVIDs, video.ID)
	}
	likesVideoInFeedListMap, err := services.GetLikesVideoIDsByUserIDAndVideoIDs(feedRequest.User.ID, feedVIDs)
	if err != nil {
		logrus.Error("get likes video map failed, err: ", err)
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrServerInternal
		p.Context.JSON(http.StatusOK, res)
		return
	}
	//debug
	for _, video := range videoModels {
		app.ZeroCheck(video.Author.ID)
	}
	res.SetValues(videoModels, likesVideoInFeedListMap, FollowedMap)
	res.CommonResponse.StatusCode = response.Success
	p.Context.JSON(http.StatusOK, res)
}

// 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func FeedVideoListHandler(c *gin.Context) {

	proxy := NewProxyFeedVideoList(c)

	var feedRequest FeedVideeDTO
	feedRequest.LatestTime = c.DefaultQuery("latest_time", fmt.Sprint(time.Now().Unix()))
	var ok bool
	user, ok := c.Get(app.UserKeyName)
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
