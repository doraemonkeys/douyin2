package main

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/Doraemonkeys/douyin2/config"
	"github.com/spf13/viper"
)

func main() {
	//os.Create("./config/conf/config.yaml")
	//生成的 example.yaml 为配置文件示例, 使用时重命名为 config.yaml
	generaetEmptyConfig()
}

func generaetEmptyConfig() {
	fileName := "example.yaml"
	targetPath := "./config/conf"

	conf := config.Config{}
	conf.Mysql.Host = "127.0.0.1"
	conf.Mysql.Port = 3306
	conf.Mysql.Username = "root"
	conf.Mysql.Password = "123456"
	conf.Mysql.Dbname = "douyin2"
	conf.Mysql.Timeout = "10s"
	conf.Mysql.MaxIdleConns = 10
	conf.Mysql.MaxOpenConns = 100
	conf.Log.Level = "trace"
	conf.Log.Path = "./log_file"
	conf.Log.PanicLogName = "gin_panic.log"
	conf.JwtSignKeyHex = "7E6F6D6E6B6C6B6A6968676665646362"
	conf.JwtSecretHex = "A1B2C3D4E5F6A1B2C3D4E5F6A1B2C3D4"
	conf.ServerPort = "6969"
	conf.Vedio.BasePath = "./uploads"
	conf.Vedio.UrlPrefix = "static"
	conf.Vedio.Domain = "http://192.168.1.105"
	err := saveNewConfig(conf, fileName, "yaml", targetPath)
	if err != nil {
		panic(err)
	}
	readConfig(filepath.Join(targetPath, fileName))
	fmt.Println(ReadTest)
}

// saveNewConfig 保存配置文件,需要传入结构体对象,配置文件名,配置文件类型,配置文件路径。
// tag必须含有mapstructure
func saveNewConfig(obj interface{}, configName, configType, configPath string) error {
	vip := viper.New()
	//反射拿到结构体的tag和value
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Tag.Get("mapstructure")
		if key == "omitempty" {
			continue
		}
		value := v.Field(i).Interface()
		//fmt.Println(key, value)
		vip.Set(key, value)
	}

	vip.SetConfigName(configName) // name of config file (without extension)
	vip.SetConfigType(configType) // REQUIRED if the config file does not have the extension in the name

	err := vip.WriteConfigAs(filepath.Join(configPath, configName))
	if err != nil {
		return err
	}
	return nil
}

var ReadTest config.Config

func readConfig(file string) {
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置文件失败, error:" + err.Error())
	}

	err = viper.Unmarshal(&ReadTest)
	if err != nil {
		panic("解析配置文件失败, error:" + err.Error())
	}
}
