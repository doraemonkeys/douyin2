package follow

import (
	"strconv"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/msgQueue"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PostFollowActionDTO struct {
	Token string `json:"token"`
	//对方用户id
	ToUserID uint `json:"to_user_id"`
	//1-关注，2-取消关注
	ActionType string `json:"action_type"`
}

const (
	ActionType_Follow   = "1"
	ActionType_Unfollow = "2"
)

const (
	PostFollowActionDTOtag_token       = "token"
	PostFollowActionDTOtag_to_user_id  = "to_user_id"
	PostFollowActionDTOtag_action_type = "action_type"
)

const (
	// 关注成功
	SuccessFollowed = "关注成功"
	// 取消关注成功
	SuccessUnfollowed = "取消关注成功"
)

func (p *PostFollowActionDTO) getAndCheckDTO(c *gin.Context) {
	// query toUserID
	var exist1 bool
	var exist2 bool
	strID, exist1 := c.GetQuery(PostFollowActionDTOtag_to_user_id)
	p.ActionType, exist2 = c.GetQuery(PostFollowActionDTOtag_action_type)
	if !exist1 || !exist2 {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	if p.ActionType != ActionType_Follow && p.ActionType != ActionType_Unfollow {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	uintID, err := strconv.ParseUint(strID, 10, 64)
	if err != nil || uintID == 0 {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	p.ToUserID = uint(uintID)
}

func PostFollowActionHandler(c *gin.Context) {
	var p PostFollowActionDTO
	p.getAndCheckDTO(c)
	msg := p.newMsg(c)
	que := msgQueue.GetFollowMQ()
	que.Push(msg)
	response.ResponseSuccess(c, SuccessFollowed)
}

func (p *PostFollowActionDTO) newMsg(c *gin.Context) msgQueue.FollowMsg {
	var msg msgQueue.FollowMsg
	msg.ToUserID = p.ToUserID
	msg.ActionType = p.ActionType
	user := c.MustGet(app.UserKeyName).(app.User)
	msg.UserID = user.ID
	return msg
}

type QueryFollowListDTO struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
}

const (
	QueryFollowListDTOtag_UserID = "user_id"
)

func (p *QueryFollowListDTO) getAndCheckDTO(c *gin.Context) {
	var exist bool
	strID, exist := c.GetQuery(QueryFollowListDTOtag_UserID)
	if !exist {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	uintID, err := strconv.ParseUint(strID, 10, 64)
	if err != nil || uintID == 0 {
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	p.UserID = uint(uintID)
}

func QueryFollowListHandler(c *gin.Context) {
	var p QueryFollowListDTO
	p.getAndCheckDTO(c)
	followList, err := services.QueryFollowUserListByUserID(p.UserID)
	if err != nil {
		logrus.Error("QueryFollowListHandler err:", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	var res response.QueryFollowListResponse
	//查询查询者是否关注了被查询者的关注列表中的用户
	user := c.MustGet(app.UserKeyName).(app.User)
	var UserList = make([]response.UserList, len(followList))

	if user.ID == p.UserID {
		// 自己查自己
		followMyself := services.QueryUserFollowed(user.ID, user.ID)
		for i, v := range followList {
			UserList[i].SetValue(v, followMyself)
			app.ZeroCheck(UserList[i].ID)
		}
		res.UserList = UserList
		res.StatusCode = response.Success
		c.JSON(200, res)
		return
	}

	queryIDs := make([]uint, len(followList))
	for i, v := range followList {
		queryIDs[i] = v.ID
	}
	followedMap, err := services.QueryUserFollowedMap(user.ID, queryIDs)
	if err != nil {
		logrus.Error("QueryFollowListHandler err:", err)
		response.ResponseError(c, response.ErrServerInternal)
		return
	}
	for i, v := range followList {
		UserList[i].SetValue(v, followedMap[v.ID])
		app.ZeroCheck(UserList[i].ID)
	}
	//debug
	for _, v := range UserList {
		logrus.Debug("queryerID: ", user.ID, " queryID: ", p.UserID, " followID: ", v.ID, " followed: ", v.IsFollow)
	}
	res.UserList = UserList
	res.StatusCode = response.Success
	c.JSON(200, res)
}
