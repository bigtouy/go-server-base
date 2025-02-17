package global

import (
	"go-server-base/configs"
	"go-server-base/init/session/psession"

	"go-server-base/init/cache/badger_db"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var (
	DB      *gorm.DB
	LOG     *logrus.Logger
	CONF    configs.ServerConfig
	Viper   *viper.Viper
	SESSION *psession.PSession
	CACHE   *badger_db.Cache
	VALID   *validator.Validate
)
