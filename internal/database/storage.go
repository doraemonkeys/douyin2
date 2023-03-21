package database

import (
	"sync"

	"github.com/Doraemonkeys/douyin2/config"
	"github.com/Doraemonkeys/douyin2/internal/pkg/storage"
)

var videoSaver *storage.LocalDouyinVedioSaver
var videoInitOnce sync.Once

func GetVideoSaver() storage.VideoStorageService[storage.SimpleObject] {
	videoInitOnce.Do(func() {
		videoSaver = storage.InitLocalOSS(GetMysqlDB(), config.GetVedioConfig().BasePath)
	})
	return videoSaver
}
