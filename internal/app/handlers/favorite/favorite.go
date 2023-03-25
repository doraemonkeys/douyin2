package favorite

import (
	"strconv"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/Doraemonkeys/douyin2/internal/msgQueue"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const FavoriteSuccess = "点赞成功"

type PostFavorDTO struct {
	VideoID int    `json:"video_id"`
	Token   string `json:"token"`
	// 1-点赞，2-取消点赞
	ActionType string `json:"action_type"`
}

const (
	PostFavorDTO_VideoID   = "video_id"
	PostFavorDTO_Token     = "token"
	PostFavorDTO_ActionTyp = "action_type"
)

func PostFavorHandler(c *gin.Context) {
	var postFavorDTO PostFavorDTO
	postFavorDTO.VideoID, _ = strconv.Atoi(c.Query(PostFavorDTO_VideoID))
	postFavorDTO.ActionType = c.Query(PostFavorDTO_ActionTyp)
	logrus.Debug("postFavorDTO: ", postFavorDTO)
	// 获取用户
	user := c.MustGet(app.UserKeyName).(app.User)
	// 点赞
	if postFavorDTO.ActionType != "1" && postFavorDTO.ActionType != "2" {
		logrus.Debug("invalid action type ", postFavorDTO.ActionType)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	action, _ := strconv.Atoi(postFavorDTO.ActionType)

	// 发送到消息队列
	mq := msgQueue.GetFavoriteMQ()
	mq.Push(msgQueue.FavoriteMSg{
		VideoID:    uint(postFavorDTO.VideoID),
		UserID:     user.ID,
		ActionType: action,
	})
	logrus.Debug("send to mq success")
	response.ResponseSuccess(c, FavoriteSuccess)
}

type QueryFavorVideoListDTO struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
}

const (
	QueryFavorVideoListDTO_UserID = "user_id"
)

func QueryFavorVideoListHandler(c *gin.Context) {
	var queryFavorVideoListDTO QueryFavorVideoListDTO
	// err := c.ShouldBind(&queryFavorVideoListDTO)
	// if err != nil {
	// 	logrus.Error("bind queryFavorVideoListDTO failed, err:", err)
	// 	response.ResponseError(c, response.ErrInvalidParams)
	// 	return
	// }
	strID, ok := c.GetQuery(QueryFavorVideoListDTO_UserID)
	UintID, err := strconv.ParseUint(strID, 10, 64)
	if !ok || err != nil {
		logrus.Error("get queryFavorVideoListDTO failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	queryFavorVideoListDTO.UserID = uint(UintID)

	// 从cache中获取用户的点赞列表
	likesCacher := database.GetUserFavoriteCacher()
	likes, eixst := likesCacher.Get(queryFavorVideoListDTO.UserID)
	if eixst {
		// cache hit
		queryFavorVideoListHandler_CacheHit(c, queryFavorVideoListDTO, likes)
		return
	}
	// cache miss
	queryFavorVideoListHandler_CacheMiss(c, queryFavorVideoListDTO)
}

func queryFavorVideoListHandler_CacheHit(c *gin.Context, queryFavorVideoListDTO QueryFavorVideoListDTO, likes models.UserLikeCacheModel) {
	var videoLikeList []models.UserLike_VideoAndAuthor
	var res response.QueryFavorVideoListResponse
	likes.VideoIDMap.Range(
		func(key, value interface{}) bool {
			videoLikeList = append(videoLikeList, value.(models.UserLike_VideoAndAuthor))
			return true
		},
	)
	videoIDs := make([]uint, len(videoLikeList))
	for i, val := range videoLikeList {
		videoIDs[i] = val.VideoID
	}
	app.ZeroListCheck(videoLikeList)
	videoAndAuthorInfos, err := services.GetVideoListAndAuthorByVideoIDList(videoIDs)
	if err != nil {
		logrus.Error("get video list failed, err:", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	var FollowedMap = make(map[uint]bool)
	user := c.MustGet(app.UserKeyName).(app.User)
	if user.ID == queryFavorVideoListDTO.UserID {
		// 自己查询自己的点赞列表，可以直接从cache中判断是否关注视频作者
		for _, val := range videoLikeList {
			FollowedMap[val.VideoID] = val.IsFollowed
		}
	} else {
		// 查询别人的点赞列表，需要从数据库中获取
		queryIdList := make([]uint, len(videoAndAuthorInfos))
		for i, val := range videoAndAuthorInfos {
			queryIdList[i] = val.AuthorID
		}
		FollowedMap, err = services.QueryFollowedMapByUserIDList(user.ID, queryIdList)
		if err != nil {
			logrus.Error("get followed map failed, err:", err)
			response.ResponseError(c, response.ErrServerInternal)
			return
		}
	}
	res.SetValues(videoAndAuthorInfos, FollowedMap)
	res.StatusCode = response.Success
	c.JSON(200, res)
}

func queryFavorVideoListHandler_CacheMiss(c *gin.Context, queryFavorVideoListDTO QueryFavorVideoListDTO) {
	var res response.QueryFavorVideoListResponse
	// cache未命中，从数据库中获取
	favorVideoIDList, err := services.QueryFavorVideoIDListByUserID(queryFavorVideoListDTO.UserID)
	if err != nil {
		logrus.Error("get favor list failed, err:", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	app.ZeroListCheck(favorVideoIDList)
	// 获取视频列表
	videoList, err := services.GetVideoListAndAuthorByVideoIDList(favorVideoIDList)
	if err != nil {
		logrus.Error("get video list failed, err:", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	// 获取视频作者是否被查询者关注
	user := c.MustGet(app.UserKeyName).(app.User)
	queryIdList := make([]uint, len(videoList))
	for i, val := range videoList {
		app.ZeroCheck(val.AuthorID)
		queryIdList[i] = val.AuthorID
	}
	followedMap, err := services.QueryFollowedMapByUserIDList(user.ID, queryIdList)
	if err != nil {
		logrus.Error("get followed map failed, err:", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	res.SetValues(videoList, followedMap)
	res.StatusCode = response.Success
	c.JSON(200, res)
}
