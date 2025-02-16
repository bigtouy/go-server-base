package global

import (
	"github.com/spf13/viper"
	"go-server-base/configs"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var (
	DB    *gorm.DB
	LOG   *logrus.Logger
	CONF  configs.ServerConfig
	Viper *viper.Viper
)
