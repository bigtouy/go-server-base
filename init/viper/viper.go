package viper

import (
	"fmt"
	"go-server-base/configs"
	"go-server-base/global"
	"os"
	"path"

	"github.com/spf13/viper"
)

func Init() {
	v := viper.NewWithOptions()
	// 从环境变量获取配置文件名
	mode := os.Getenv("MODE")
	fmt.Println("MODE:", mode)

	var configName = "app."

	switch mode {
	case "dev":
		configName = configName + "dev"
	case "prod":
		configName = configName + "prod"
	default:
		configName = configName + "prod"
	}

	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}

	serverConfig := configs.ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	global.CONF = serverConfig
	global.CONF.Log.LogPath = path.Join(global.CONF.System.BaseDir, global.CONF.App.Name, global.CONF.Log.LogPath)
	global.CONF.Sqlite.DbPath = path.Join(global.CONF.System.BaseDir, global.CONF.App.Name, global.CONF.Sqlite.DbPath)
	global.Viper = v
}
