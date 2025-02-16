package viper

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"go-server-base/cmd/server/conf"
	"go-server-base/configs"
	"go-server-base/global"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

func Init() {
	v := viper.NewWithOptions()

	//v.AddConfigPath(path.Join("/opt/go-server-base/conf")) // search config in the working directory
	//v.SetConfigName("app")                                      // name of config file (without extension)
	v.SetConfigType("yaml")
	config := configs.ServerConfig{}
	// 从环境变量获取配置文件名
	mode := os.Getenv("MODE")

	var configYaml []byte

	switch mode {
	case "dev":
		configYaml = conf.AppDevYaml
	case "prod":
		configYaml = conf.AppProdYaml
	default:
		configYaml = conf.AppDevYaml
	}

	if err := yaml.Unmarshal(configYaml, &config); err != nil {
		panic(err)
	}
	v.SetConfigName("config")

	// 获取当前工作目录
	p, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前路径:", err)
		return
	}
	v.AddConfigPath(".")

	// 打印当前路径
	fmt.Println("当前路径:", p)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	reader := bytes.NewReader(configYaml)
	if err := v.ReadConfig(reader); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
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
