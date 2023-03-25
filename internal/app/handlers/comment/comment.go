package comment

import (
	"net/http"
	"strconv"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/sirupsen/logrus"

	"github.com/Doraemonkeys/douyin2/internal/msgQueue"
	"github.com/gin-gonic/gin"
)

type PostCommentDTO struct {
	VideoID    string `json:"video_id"`
	Token      string `json:"token"`
	ActionType string `json:"action_type"`
	//用户填写的评论内容，在action_type=1的时候使用
	CommentText string `json:"comment_text"`
	//要删除的评论id，在action_type=2的时候使用
	CommentId string `json:"comment_id"`
}

const (
	PostCommentDTO_VideoID     = "video_id"
	PostCommentDTO_ActionType  = "action_type"
	PostCommentDTO_CommentId   = "comment_id"
	PostCommentDTO_CommentText = "comment_text"
)

const (
	PostCommentDTO_ActionType_Add    = "1"
	PostCommentDTO_ActionType_Delete = "2"
)

const (
	SuccComment = "评论成功"
	SuccDelete  = "删除成功"
)

// 应该确保dto中的字段都是合法的
func (dto *PostCommentDTO) newMsg(c *gin.Context) msgQueue.CommentMsg {
	var Msg msgQueue.CommentMsg
	uintID, _ := strconv.ParseUint(dto.VideoID, 10, 64)
	Msg.VideoID = uint(uintID)
	uintCommentID, _ := strconv.ParseUint(dto.CommentId, 10, 64)
	Msg.CommentId = uint(uintCommentID)
	Msg.ActionType = dto.ActionType
	Msg.CommentText = dto.CommentText
	user := c.MustGet(app.UserKeyName).(app.User)
	Msg.CommenterID = user.ID
	return Msg
}

func PostCommentHandler(c *gin.Context) {
	var dto PostCommentDTO
	dto.getAndCheckPostCommentDTO(c)
	mq := msgQueue.GetCommentMQ()
	var Msg msgQueue.CommentMsg = dto.newMsg(c)
	mq.Push(Msg)
	response.ResponseSuccess(c, SuccComment)
}

func (dto *PostCommentDTO) getAndCheckPostCommentDTO(c *gin.Context) {
	var exist1 bool
	var exist2 bool
	dto.VideoID, exist1 = c.GetQuery(PostCommentDTO_VideoID)
	dto.ActionType, exist2 = c.GetQuery(PostCommentDTO_ActionType)
	if !exist1 || !exist2 {
		response.ResponseError(c, response.ErrInvalidParams)
	}
	var exist bool
	if dto.ActionType == PostCommentDTO_ActionType_Add {
		dto.CommentText, exist = c.GetQuery(PostCommentDTO_CommentText)
		if !exist {
			response.ResponseError(c, response.ErrInvalidParams)
			return
		}
	} else if dto.ActionType == PostCommentDTO_ActionType_Delete {
		dto.CommentId, exist = c.GetQuery(PostCommentDTO_CommentId)
		if !exist {
			response.ResponseError(c, response.ErrInvalidParams)
			return
		}
	} else {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
}

type QueryCommentListDTO struct {
	VideoID uint   `json:"video_id"`
	Token   string `json:"token"`
}

const (
	QueryCommentListDTO_VideoID = "video_id"
)

const CommentResFormat = "2006-01-02 15:04:05"

func QueryCommentListHandler(c *gin.Context) {
	var dto QueryCommentListDTO
	dto.getAndCheckQueryCommentListDTO(c)

	// get comments from cache
	commentCache := database.GetVideoCommentCacher()
	commentsCache, exist := commentCache.Get(dto.VideoID)
	if exist {
		logrus.Info("QueryCommentListHandler: cache hit ", commentsCache.CacheMap)
		queryCommentListHandler_CacheHit(c, commentsCache)
		return
	}
	// get comments from db
	queryCommentListHandler_CacheMiss(c, dto)
}

func queryCommentListHandler_CacheHit(c *gin.Context, commentsCache models.CommentCacheModel) {
	var res response.QueryCommentListResponse
	var CommentList = make([]response.CommentList, len(commentsCache.CacheMap))
	var commenterIDMap = make(map[uint]struct{}, len(commentsCache.CacheMap))
	commentsCache.MapLock.RLock()
	var i = 0
	for _, comment := range commentsCache.CacheMap {
		CommentList[i].ID = int(comment.ID)
		CommentList[i].Content = comment.Content
		date := comment.CreatedAt.Format(CommentResFormat)
		CommentList[i].CreateDate = date
		if comment.Commenter.ID == 0 {
			logrus.Error("commenterID is 0, comment: ", comment.Content, " commentID: ", comment.ID)
		}
		CommentList[i].User.ID = int(comment.UserID)
		commenterIDMap[comment.UserID] = struct{}{}
		logrus.Debug("commenterID: ", comment.UserID, " comment: ", comment.Content)
		i++
	}
	commentsCache.MapLock.RUnlock()

	commenterMap, err := services.GetUserMapByUserIdMap(commenterIDMap)
	if err != nil {
		logrus.Error("QueryCommentListHandler: services.GetUserMapByUserIdMap error: ", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	//logrus.Debug("QueryCommentListHandler: commenterMap: ", commenterMap)
	queryer := c.MustGet(app.UserKeyName).(app.User)
	followedMap, err := services.QueryFollowedMapByUserIDMap(queryer.ID, commenterIDMap)
	if err != nil {
		logrus.Error("QueryCommentListHandler: services.QueryFollowedMapByUserIDMap error: ", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	var temp response.User
	for k, comment := range CommentList {
		temp.SetValue(commenterMap[uint(comment.User.ID)], followedMap[uint(comment.User.ID)])
		logrus.Debug("QueryCommentListHandler: temp: ", temp)
		CommentList[k].User = temp
	}
	res.CommentList = CommentList
	res.StatusCode = response.Success
	for _, val := range res.CommentList {
		app.ZeroCheck(val.ID, val.User.ID)
	}
	c.JSON(http.StatusOK, res)
}

func queryCommentListHandler_CacheMiss(c *gin.Context, dto QueryCommentListDTO) {
	var res response.QueryCommentListResponse
	comments, err := services.QueryCommentListWithCommenterByVideoID(dto.VideoID)
	if err != nil {
		logrus.Error("QueryCommentListHandler: services.QueryCommentListWithCommenterByVideoID error: ", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	commenterIDMap := make(map[uint]struct{}, len(comments))
	for _, comment := range comments {
		commenterIDMap[comment.UserID] = struct{}{}
	}
	queryer := c.MustGet(app.UserKeyName).(app.User)
	followedMap, err := services.QueryFollowedMapByUserIDMap(queryer.ID, commenterIDMap)
	if err != nil {
		logrus.Error("QueryCommentListHandler: services.QueryFollowedMapByUserIDMap error: ", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	var CommentList = make([]response.CommentList, len(comments))
	var temp response.User
	for k, comment := range comments {
		CommentList[k].ID = int(comment.ID)
		CommentList[k].Content = comment.Content
		date := comment.CreatedAt.Format(CommentResFormat)
		CommentList[k].CreateDate = date
		temp.SetValue(comment.Commenter, followedMap[comment.UserID])
		CommentList[k].User = temp
	}
	res.CommentList = CommentList
	res.StatusCode = response.Success
	for _, val := range res.CommentList {
		app.ZeroCheck(val.ID, val.User.ID)
	}
	c.JSON(http.StatusOK, res)
}

func (dto *QueryCommentListDTO) getAndCheckQueryCommentListDTO(c *gin.Context) {
	strVideoID, exist := c.GetQuery(QueryCommentListDTO_VideoID)
	if !exist {
		response.ResponseError(c, response.ErrInvalidParams)
	}
	uintVideoID, err := strconv.ParseUint(strVideoID, 10, 64)
	if err != nil {
		logrus.Error("QueryCommentListHandler: strconv.ParseUint error: ", err)
		response.ResponseError(c, response.ErrInvalidParams)
	}
	dto.VideoID = uint(uintVideoID)
}
