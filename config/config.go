package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Doraemonkeys/douyin2/pkg/log"
	"github.com/Doraemonkeys/douyin2/utils"
	"github.com/spf13/viper"
)

var allConfig Config

func init() {
	file := "config.yaml"
	basePath := "./config/conf"
	if !utils.FileOrDirIsExist(filepath.Join(basePath, file)) {
		os.Create(filepath.Join(basePath, file))
		fmt.Println("配置文件不存在,已创建配置文件,请修改配置文件后重启程序")
	}
	readConfig(filepath.Join(basePath, file))
	if allConfig.Vedio.BasePath == "" {
		panic("请配置视频存储路径")
	}
	if allConfig.Vedio.UrlPrefix == "" {
		panic("请配置视频访问路径")
	}
}

func GetMysqlConfig() MysqlConfig {
	return allConfig.Mysql
}

func GetGlobalLoggerConfig() log.LogConfig {
	return defaultLogConfig()
}

func GetLogConfig() LogConfig {
	return allConfig.Log
}

func GetJwtConfig() JwtConfig {
	return JwtConfig{allConfig.JwtSignKeyHex, allConfig.JwtSecretHex}
}

func GetServerPort() string {
	return allConfig.ServerPort
}

func GetVedioConfig() VedioConfig {
	return allConfig.Vedio
}

func defaultLogConfig() log.LogConfig {
	var Config log.LogConfig
	Config.DateSplit = true
	Config.LogPath = allConfig.Log.Path
	Config.LogLevel = allConfig.Log.Level
	Config.ShowShortFileInConsole = true
	return Config
}

func readConfig(file string) {
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置文件失败, error:" + err.Error())
	}

	err = viper.Unmarshal(&allConfig)
	if err != nil {
		panic("解析配置文件失败, error:" + err.Error())
	}
}
