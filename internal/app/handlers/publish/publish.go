package publish

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/Doraemonkeys/douyin2/internal/app"
	"github.com/Doraemonkeys/douyin2/internal/app/handlers/response"
	"github.com/Doraemonkeys/douyin2/internal/app/models"
	"github.com/Doraemonkeys/douyin2/internal/app/services"
	"github.com/Doraemonkeys/douyin2/internal/database"
	"github.com/Doraemonkeys/douyin2/internal/pkg/storage"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	//VideoTitle is Empty
	ErrorVideoTitleEmpty  = "视频标题为空"
	ErrVedioAlreadyExists = "视频已经存在"

	//Success
	SuccessVedioUpload = "视频上传成功"
)

type PublishVedioDTO struct {
	Data  []byte `json:"data"`
	Token string `json:"token"`
	Title string `json:"title"`
}

func titleCheck(title string) bool {
	if utils.StrLen(title) > models.VideoTitleMaxRuneLength {
		return false
	}
	if len(title) > models.VideoTitleMaxByteLength {
		return false
	}
	if len(title) == 0 {
		return false
	}
	return true
}

func PublishVedioHandler(c *gin.Context) {
	var publishRequest PublishVedioDTO
	FileHeader, err := c.FormFile("data")
	if err != nil {
		logrus.Error("receive file failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	//title check
	publishRequest.Title = c.PostForm("title")
	if publishRequest.Title == "" {
		response.ResponseError(c, ErrorVideoTitleEmpty)
	}
	if !titleCheck(publishRequest.Title) {
		logrus.Trace("title is too long", publishRequest.Title, "len:", len(publishRequest.Title))
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	err = publishRequest.ReceiveFile(FileHeader)
	if err != nil {
		logrus.Error("receive file failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	// 保存视频
	var videoObject storage.SimpleObject
	videoObject.Name = FileHeader.Filename
	videoObject.Data = publishRequest.Data
	storageID, err := database.GetVideoSaver().SaveUnique(videoObject)
	if errors.Is(err, storage.ErrVedioAlreadyExists) {
		response.ResponseError(c, ErrVedioAlreadyExists)
		return
	}
	if err != nil {
		logrus.Error("save vedio failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	// 返回success
	response.ResponseSuccess(c, SuccessVedioUpload)

	videoUrl, CoverUrl, err := database.GetVideoSaver().GetURL(storageID)
	if err != nil {
		logrus.Error("get vedio url failed, err:", err)
		return
	}
	logrus.Trace("videoUrl:", videoUrl, "CoverUrl:", CoverUrl)
	//获取User
	user := c.MustGet("user").(app.User)
	var videoModel models.VideoModel
	videoModel.Title = publishRequest.Title
	videoModel.StorageID = storageID
	videoModel.AuthorID = user.ID
	videoModel.URL = videoUrl
	videoModel.CoverURL = CoverUrl

	//保存视频信息到数据库
	err = services.CreateVedio(&videoModel)
	if err != nil {
		logrus.Error("save vedio info failed, err:", err)
		return
	}
}

func (p *PublishVedioDTO) ReceiveFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	p.Data = make([]byte, fileHeader.Size)
	_, err = file.Read(p.Data)
	if err != nil {
		return err
	}
	return nil
}

func (p *PublishVedioDTO) SaveFile(newFilePath string) error {
	file, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(p.Data)
	return err
}

// 用户的视频发布列表，直接列出用户所有投稿过的视频
type QueryVideoListDTO struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

func QueryPublishListHandler(c *gin.Context) {
	var queryVideoListDTO QueryVideoListDTO
	err := c.ShouldBind(&queryVideoListDTO)
	if err != nil {
		logrus.Error("bind queryVideoListDTO failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}

	// 获取被查询的用户
	targetUser, err := services.GetUserById(uint(queryVideoListDTO.UserId))
	if err != nil && err.Error() == response.ErrUserNotExists {
		response.ResponseError(c, response.ErrUserNotExists)
		return
	}
	if err != nil {
		logrus.Error("get user failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	// 获取被查询的用户发布列表
	targetUserVideoPublishList, err := services.QueryPublishListByAuthorID(uint(queryVideoListDTO.UserId))
	if err != nil {
		logrus.Error("get video list failed, err:", err)
		response.ResponseError(c, response.ErrInvalidParams)
		return
	}
	if len(targetUserVideoPublishList) == 0 {
		logrus.Trace("video list is empty")
		response.ResponseSuccess(c, response.EmptyVideoList)
		return
	}
	// 获取查询者
	user := c.MustGet("user").(app.User)
	// 查询查询者所有喜欢的视频
	var likeVideoIDs []uint
	// 1. from cache
	cacher := database.GetUserFavoriteCacher()
	cacheData, exist := cacher.Get(user.ID)
	if exist {
		cacheData.VideoIDMap.Range(func(key, value interface{}) bool {
			likeVideoIDs = append(likeVideoIDs, key.(models.UserLike_VideoAndAuthor).VideoID)
			return true
		})
	} else {
		// 2. from db
		likeVideos, _, _ := services.QueryLikeVideoListByUserID(user.ID)
		for _, video := range likeVideos {
			likeVideoIDs = append(likeVideoIDs, video.ID)
		}
	}
	var likeVideoListMap = make(map[uint]bool, len(likeVideoIDs))
	for _, id := range likeVideoIDs {
		likeVideoListMap[id] = true
	}
	// 查询者是否关注了被查询的用户
	isFollowed := services.QueryUserFollowed(user.ID, targetUser.ID)
	// 返回视频列表
	var res response.QueryVideoListResponse
	res.SetValues(targetUserVideoPublishList, targetUser, likeVideoListMap, isFollowed)

	res.StatusCode = response.Success
	res.StatusMsg = response.QuerySuccessMsg
	c.JSON(http.StatusOK, res)

}
