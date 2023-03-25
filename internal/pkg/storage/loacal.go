package storage

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/Doraemonkeys/douyin2/config"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SimpleObject struct {
	// Object name
	Name string `json:"name"`
	// Data
	Data []byte `json:"data"`
}

type VedioObjectModel struct {
	UID        uint   `json:"uid" gorm:"primary_key"`      // 视频uid
	Path       string `json:"path" gorm:"size:200"`        // 视频路径
	SHA1       string `json:"sha1" gorm:"size:100,index"`  // 视频sha1
	CoverPath  string `json:"cover_path" gorm:"size:200"`  // 封面路径
	OriginName string `json:"origin_name" gorm:"size:100"` // 原始文件名
}

type LocalDouyinVedioSaver struct {
	// Base path
	BasePath string
	db       *gorm.DB
}

var (
	ErrVedioAlreadyExists = errors.New("视频已经存在")
)

var globalSaver *LocalDouyinVedioSaver
var once sync.Once

func InitLocalOSS(mysql *gorm.DB, basePath string) *LocalDouyinVedioSaver {
	once.Do(func() {
		globalSaver = &LocalDouyinVedioSaver{
			BasePath: basePath,
			db:       mysql,
		}
		err := globalSaver.db.AutoMigrate(&VedioObjectModel{})
		if err != nil {
			panic("初始化视频存储器失败")
		}
	})
	return globalSaver
}

// Save object to local,Return VedioUID, error
func (s *LocalDouyinVedioSaver) Save(obj SimpleObject) (uid uint, err error) {
	panic("implement me")
}

func (s *LocalDouyinVedioSaver) Delete(uid uint) error {
	panic("implement me")
}

func (s *LocalDouyinVedioSaver) Get(uid uint) (SimpleObject, error) {
	panic("implement me")
}

// QueryOrignalName 根据视频uid查询视频原始文件名
func (s *LocalDouyinVedioSaver) QueryOrignalName(uid uint) (string, error) {
	var video VedioObjectModel
	err := s.db.Where("uid = ?", uid).Take(&video).Error
	if err != nil {
		return "", err
	}
	return video.OriginName, nil
}

func (s *LocalDouyinVedioSaver) SaveAndQueryExist(obj SimpleObject) (uid uint, exist bool, err error) {
	uid, err = s.SaveUnique(obj)
	if err != nil {
		return uid, false, nil
	}
	if errors.Is(err, ErrVedioAlreadyExists) {
		return uid, true, nil
	}
	return 0, false, err
}

func (s *LocalDouyinVedioSaver) QueryExistBySHA1(sha1 string) bool {
	sha1 = strings.ToLower(sha1)
	var video VedioObjectModel
	err := s.db.Where("sha1 = ?", sha1).Find(&video).Error
	if err == gorm.ErrRecordNotFound || video.UID == 0 {
		return false
	}
	return true
}

// 返回视频uid,如果视频已经存在则返回错误ErrVedioAlreadyExists
func (s *LocalDouyinVedioSaver) SaveUnique(video SimpleObject) (uid uint, err error) {
	// 计算视频sha1
	sha1, err := utils.GetSha1(video.Data)
	if err != nil {
		return
	}
	// 查询视频是否已经存在
	if s.QueryExistBySHA1(sha1) {
		return 0, ErrVedioAlreadyExists
	}

	newVideoFileName := generateNewVideoName()
	// 如果视频文件名有后缀则保留后缀
	if filepath.Ext(video.Name) != "" {
		newVideoFileName = newVideoFileName + filepath.Ext(video.Name)
	}
	videoPath, coverPath := generateNewVedioAndCoverPath(s.BasePath, newVideoFileName)
	if !utils.FileOrDirIsExist(filepath.Dir(videoPath)) {
		if err := utils.CreateDir(filepath.Dir(videoPath)); err != nil {
			return 0, err
		}
	}
	//保存文件
	err = video.saveFile(videoPath)
	if err != nil {
		return
	}
	//保存封面
	err = utils.CreateCoverFromLocal(videoPath, coverPath)
	if err != nil {
		return
	}
	//保存数据库
	vedioObjectModel := VedioObjectModel{
		Path:       videoPath,
		SHA1:       sha1,
		CoverPath:  coverPath,
		OriginName: video.Name,
	}
	err = s.CreateVedio(&vedioObjectModel)
	if err != nil {
		return
	}
	return vedioObjectModel.UID, nil
}

func (s *LocalDouyinVedioSaver) GetURL(uid uint) (videoUrl string, coverUrl string, err error) {
	var video VedioObjectModel
	err = s.db.Where("uid = ?", uid).Take(&video).Error
	if err != nil {
		return
	}
	// windows下路径使用\分割，需要替换为/
	video.Path = strings.ReplaceAll(video.Path, "\\", "/")
	video.CoverPath = strings.ReplaceAll(video.CoverPath, "\\", "/")

	// 去除路径前面的/
	if video.Path[0] == '/' {
		video.Path = video.Path[1:]
	}
	if video.CoverPath[0] == '/' {
		video.CoverPath = video.CoverPath[1:]
	}

	// 去除basePath前后的/
	// ./upload --> upload
	// ./upload/ --> upload
	// upload/ --> upload
	// upload --> upload
	// video/upload --> video/upload
	// ./video/upload --> video/upload
	tempBasePath := strings.ReplaceAll(s.BasePath, "\\", "/")
	if tempBasePath != "" && tempBasePath[0] == '.' {
		tempBasePath = tempBasePath[1:]
	}
	if tempBasePath != "" && tempBasePath[len(tempBasePath)-1] == '/' {
		tempBasePath = tempBasePath[:len(tempBasePath)-1]
	}
	if tempBasePath != "" && tempBasePath[0] == '/' {
		tempBasePath = tempBasePath[1:]
	}

	// domain:port/urlPrefix/xxx/xxx/xxx.mp4
	baseUrl := config.GetVedioConfig().Domain + ":" + config.GetServerPort() + "/" + config.GetVedioConfig().UrlPrefix

	// basePath/2023/03/20/xxx.mp4 -> 2023/03/20/xxx.mp4
	subVideoPath := strings.Replace(video.Path, tempBasePath, "", 1)
	subCoverPath := strings.Replace(video.CoverPath, tempBasePath, "", 1)
	if subVideoPath[0] == '/' {
		subVideoPath = subVideoPath[1:]
	}
	if subCoverPath[0] == '/' {
		subCoverPath = subCoverPath[1:]
	}
	videoUrl = baseUrl + "/" + subVideoPath
	coverUrl = baseUrl + "/" + subCoverPath
	return
}

func (s *SimpleObject) saveFile(newFilePath string) error {
	file, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(s.Data)
	return err
}

// generateNewVedioAndCoverPath 返回保存视频文件的路径和封面文件的路径
//
// path: /data/vedio/2020/01/01/xxx.mp4
func generateNewVedioAndCoverPath(basePath string, baseName string) (string, string) {
	logrus.Trace("vedio basePath:", basePath)
	logrus.Trace("vedio baseName:", baseName)

	// 生成视频文件的路径
	vedioPath := filepath.Join(
		basePath,
		time.Now().Format("2006"),
		time.Now().Format("01"),
		time.Now().Format("02"),
		baseName,
	)
	logrus.Trace("vedioPath:", vedioPath)
	// 生成封面文件的路径
	coverPath := vedioPath + ".jpg"
	logrus.Trace("coverPath:", coverPath)
	return vedioPath, coverPath
}

// 生成新的视频文件名(时间戳+随机数,不带后缀)
func generateNewVideoName() string {
	randNum := rand.Intn(1000)
	newName := fmt.Sprintf("%d", time.Now().UnixNano()) + fmt.Sprintf("%d", randNum)
	return newName
}

func (s *LocalDouyinVedioSaver) CreateVedio(video *VedioObjectModel) error {
	return s.db.Create(video).Error
}
